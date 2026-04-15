"""
Miora AI — AgentKit Sidecar Service

A lightweight FastAPI service that wraps Coinbase AgentKit for autonomous trading
on Base Sepolia. This sidecar is called by the Go backend via HTTP.

Architecture:
    Go backend (agent_loop.go)
        → HTTP call to this sidecar (localhost:8090)
            → AgentKit SDK (coinbase-agentkit)
                → Agentic Wallet on Base Sepolia
                    → On-chain transaction

Why a sidecar?
    AgentKit SDK only exists in Python and TypeScript — not Go.
    So we run this as a separate process and the Go backend calls it via HTTP.

Endpoints:
    GET  /health   — Check if sidecar is running and agent is initialized
    GET  /wallet   — Get the Agentic Wallet address and balance
    POST /swap     — Execute a token swap using the Agentic Wallet
    POST /transfer — Transfer native ETH from the Agentic Wallet
    GET  /actions  — List all available AgentKit actions

Agentic Wallet:
    - Created and managed by Coinbase AgentKit (CDP API)
    - Private key is held by Coinbase, NOT on our server
    - User deposits funds to this wallet, agent trades with it
    - User can withdraw anytime

Runs on port 8090 by default (configurable via AGENT_PORT env var).
Start with: python main.py  or  make run-agent
"""

import os
import json
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

from coinbase_agentkit import (
    AgentKit,
    AgentKitConfig,
    CdpEvmWalletProvider,
    CdpEvmWalletProviderConfig,
    wallet_action_provider,
    erc20_action_provider,
)

load_dotenv()

# --- Global state ---
# These are initialized once on startup and reused for all requests.
wallet_provider = None  # CdpEvmWalletProvider — manages the Agentic Wallet
agent_kit = None        # AgentKit instance — provides actions (transfer, swap, etc.)


def init_agent():
    """Initialize the AgentKit wallet provider and agent instance.
    
    This creates (or restores) an Agentic Wallet on Base Sepolia using
    Coinbase Developer Platform (CDP) credentials. The wallet is managed
    by Coinbase — we never touch the private key.
    
    Required env vars:
        CDP_API_KEY_ID     — CDP API key ID (from https://portal.cdp.coinbase.com)
        CDP_API_KEY_SECRET — CDP API key secret
        CDP_NETWORK_ID     — Network to use (default: base-sepolia)
    """
    global wallet_provider, agent_kit

    api_key_id = os.getenv("CDP_API_KEY_ID")
    api_key_secret = os.getenv("CDP_API_KEY_SECRET")
    network_id = os.getenv("CDP_NETWORK_ID", "base-sepolia")

    if not api_key_id or not api_key_secret:
        print("[WARN] CDP credentials not set — agent will not be functional")
        return

    # Create wallet provider — this creates a new Agentic Wallet on first run,
    # or restores the existing one on subsequent runs.
    wallet_provider = CdpEvmWalletProvider(CdpEvmWalletProviderConfig(
        api_key_id=api_key_id,
        api_key_secret=api_key_secret,
        network_id=network_id,
    ))

    # Create AgentKit instance with wallet and action providers.
    # wallet_action_provider: get_balance, native_transfer, get_wallet_details
    # erc20_action_provider: transfer ERC-20 tokens
    agent_kit = AgentKit(AgentKitConfig(
        wallet_provider=wallet_provider,
        action_providers=[
            wallet_action_provider(),
            erc20_action_provider(),
        ],
    ))

    print(f"[AgentKit] Initialized on {network_id}")
    print(f"[AgentKit] Wallet address: {wallet_provider.get_address()}")


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Initialize agent on startup, cleanup on shutdown."""
    init_agent()
    yield


app = FastAPI(
    title="Miora AgentKit Sidecar",
    description="Coinbase AgentKit wrapper for autonomous trading on Base",
    version="1.0.0",
    lifespan=lifespan,
)


# --- Request/Response models ---

class TransferRequest(BaseModel):
    """Request to transfer native ETH from the Agentic Wallet."""
    to_address: str  # Destination wallet address
    amount_eth: str  # Amount of ETH to send (e.g. "0.001")


class SwapRequest(BaseModel):
    """Request to execute a token swap via the Agentic Wallet.
    
    For hackathon PoC, this uses native_transfer to simulate a swap.
    In production, this would use a DEX aggregator action (Uniswap/1inch).
    """
    token_address: str      # Token contract address to buy
    amount_eth: str         # Amount of ETH to spend (e.g. "0.001")
    token_symbol: str = ""  # Optional token symbol for logging


# --- Endpoints ---

@app.get("/health")
async def health():
    """Health check — returns whether the agent is initialized and wallet address."""
    return {
        "status": "ok",
        "agent_ready": agent_kit is not None,
        "wallet": wallet_provider.get_address() if wallet_provider else None,
    }


@app.get("/wallet")
async def get_wallet():
    """Get the Agentic Wallet address, network, and balance.
    
    The Agentic Wallet is the wallet that the agent uses to trade.
    User deposits funds here, agent manages them autonomously.
    """
    if not wallet_provider:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    address = wallet_provider.get_address()
    network = os.getenv("CDP_NETWORK_ID", "base-sepolia")

    # Get balance using AgentKit's built-in action
    actions = agent_kit.get_actions()
    balance_action = next((a for a in actions if a.name == "get_balance"), None)

    balance = "unknown"
    if balance_action:
        try:
            result = balance_action.invoke({})
            balance = result
        except Exception as e:
            balance = f"error: {e}"

    return {
        "address": address,
        "network": network,
        "balance": balance,
    }


@app.post("/swap")
async def execute_swap(req: SwapRequest):
    """Execute a token swap using the Agentic Wallet.
    
    Called by Go backend's agent_loop.go when:
    1. A top-scored wallet makes a trade
    2. All conditions pass (liquidity, mcap, pair age, budget, AI risk)
    3. Agent status is "active"
    
    For hackathon PoC: uses native_transfer to simulate a swap.
    In production: would integrate with Uniswap/1inch action provider
    to actually swap ETH for the target token.
    """
    if not agent_kit:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    try:
        actions = agent_kit.get_actions()
        transfer_action = next((a for a in actions if a.name == "native_transfer"), None)

        if not transfer_action:
            raise HTTPException(status_code=500, detail="Transfer action not available")

        # Execute the swap (PoC: native transfer to token address)
        result = transfer_action.invoke({
            "to": req.token_address,
            "value": req.amount_eth,
        })

        return {
            "status": "success",
            "token_address": req.token_address,
            "token_symbol": req.token_symbol,
            "amount_eth": req.amount_eth,
            "result": result,
            "agent_wallet": wallet_provider.get_address(),
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/transfer")
async def transfer(req: TransferRequest):
    """Transfer native ETH from the Agentic Wallet to any address.
    
    Used for withdrawals — user can pull funds back from the agent wallet.
    """
    if not agent_kit:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    actions = agent_kit.get_actions()
    transfer_action = next((a for a in actions if a.name == "native_transfer"), None)

    if not transfer_action:
        raise HTTPException(status_code=500, detail="Transfer action not available")

    try:
        result = transfer_action.invoke({
            "to": req.to_address,
            "value": req.amount_eth,
        })
        return {"status": "success", "result": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/actions")
async def list_actions():
    """List all available AgentKit actions.
    
    Useful for debugging — shows what the agent can do.
    """
    if not agent_kit:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    actions = agent_kit.get_actions()
    return {
        "actions": [
            {"name": a.name, "description": a.description}
            for a in actions
        ]
    }


if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("AGENT_PORT", "8090"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)

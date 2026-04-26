"""
Miora AI — AgentKit Sidecar Service

A lightweight FastAPI service that wraps Coinbase AgentKit for autonomous trading
on Base Sepolia. This sidecar is called by the Go backend via HTTP.

Endpoints:
    GET  /health   — Check if sidecar is running and agent is initialized
    GET  /wallet   — Get the Agentic Wallet address and balance
    POST /swap     — Execute a token swap using the Agentic Wallet
    POST /transfer — Transfer native ETH from the Agentic Wallet
    GET  /actions  — List all available AgentKit actions

Runs on port 8090 by default (configurable via AGENT_PORT env var).
"""

import os
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

load_dotenv()

# --- Global state ---
wallet_provider = None
agent_kit = None


def init_agent():
    """Initialize the AgentKit wallet provider using CDP Server Wallet.

    Required env vars:
        CDP_API_KEY_ID     — CDP API key ID
        CDP_API_KEY_SECRET — CDP API key secret
        CDP_WALLET_SECRET  — CDP wallet secret (for signing transactions)
        CDP_NETWORK_ID     — Network (default: base-sepolia)
    """
    global wallet_provider, agent_kit

    api_key_id = os.getenv("CDP_API_KEY_ID")
    api_key_secret = os.getenv("CDP_API_KEY_SECRET")
    wallet_secret = os.getenv("CDP_WALLET_SECRET")
    network_id = os.getenv("CDP_NETWORK_ID", "base-sepolia")

    if not api_key_id or not api_key_secret or not wallet_secret:
        missing = []
        if not api_key_id: missing.append("CDP_API_KEY_ID")
        if not api_key_secret: missing.append("CDP_API_KEY_SECRET")
        if not wallet_secret: missing.append("CDP_WALLET_SECRET")
        print(f"[WARN] Missing: {', '.join(missing)} — agent in demo mode")
        return

    try:
        from coinbase_agentkit import (
            AgentKit,
            AgentKitConfig,
            CdpEvmWalletProvider,
            CdpEvmWalletProviderConfig,
            wallet_action_provider,
            erc20_action_provider,
        )

        wallet_provider = CdpEvmWalletProvider(CdpEvmWalletProviderConfig(
            api_key_id=api_key_id,
            api_key_secret=api_key_secret,
            wallet_secret=wallet_secret,
            network_id=network_id,
        ))

        agent_kit = AgentKit(AgentKitConfig(
            wallet_provider=wallet_provider,
            action_providers=[
                wallet_action_provider(),
                erc20_action_provider(),
            ],
        ))

        print(f"[AgentKit] Initialized on {network_id}")
        print(f"[AgentKit] Wallet address: {wallet_provider.get_address()}")

    except Exception as e:
        print(f"[ERROR] Failed to initialize AgentKit: {e}")
        print("[WARN] Agent will run in demo mode")


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Init agent in lifespan. Run in thread to avoid asyncio.run() conflict."""
    if agent_kit is None:
        import asyncio
        loop = asyncio.get_event_loop()
        await loop.run_in_executor(None, init_agent)
    yield


app = FastAPI(
    title="Miora AgentKit Sidecar",
    description="Coinbase AgentKit wrapper for autonomous trading on Base",
    version="1.0.0",
    lifespan=lifespan,
)


# --- Request/Response models ---

class TransferRequest(BaseModel):
    to_address: str
    amount_eth: str


class SwapRequest(BaseModel):
    token_address: str
    amount: str
    token_symbol: str = ""
    direction: str = "buy"
    musdt_address: str = ""


# --- Endpoints ---

@app.get("/health")
async def health():
    return {
        "status": "ok",
        "agent_ready": agent_kit is not None,
        "wallet": wallet_provider.get_address() if wallet_provider else None,
    }


@app.get("/wallet")
async def get_wallet():
    if not wallet_provider:
        return {"address": "", "network": "base-sepolia", "balance": "0", "demo_mode": True}

    address = wallet_provider.get_address()
    actions = agent_kit.get_actions()
    balance_action = next((a for a in actions if a.name == "get_balance"), None)

    balance = "unknown"
    if balance_action:
        try:
            result = balance_action.invoke({})
            balance = result
        except Exception as e:
            balance = f"error: {e}"

    return {"address": address, "network": os.getenv("CDP_NETWORK_ID", "base-sepolia"), "balance": balance}


@app.post("/swap")
async def execute_swap(req: SwapRequest):
    if not agent_kit:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    try:
        actions = agent_kit.get_actions()
        transfer_action = next((a for a in actions if a.name == "transfer"), None)
        if not transfer_action:
            transfer_action = next((a for a in actions if a.name == "native_transfer"), None)
        if not transfer_action:
            raise HTTPException(status_code=500, detail="Transfer action not available")

        result = transfer_action.invoke({"to": req.token_address, "value": req.amount})

        return {
            "status": "success",
            "direction": req.direction,
            "token_address": req.token_address,
            "token_symbol": req.token_symbol,
            "amount": req.amount,
            "result": result,
            "agent_wallet": wallet_provider.get_address(),
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/transfer")
async def transfer(req: TransferRequest):
    if not agent_kit:
        raise HTTPException(status_code=503, detail="Agent not initialized")

    actions = agent_kit.get_actions()
    transfer_action = next((a for a in actions if a.name == "native_transfer"), None)
    if not transfer_action:
        raise HTTPException(status_code=500, detail="Transfer action not available")

    try:
        result = transfer_action.invoke({"to": req.to_address, "value": req.amount_eth})
        return {"status": "success", "result": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/actions")
async def list_actions():
    if not agent_kit:
        return {"actions": [], "demo_mode": True}

    actions = agent_kit.get_actions()
    return {"actions": [{"name": a.name, "description": a.description} for a in actions]}


if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("AGENT_PORT", "8090"))
    # reload=False to avoid reloader process conflict with CdpEvmWalletProvider
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=False)

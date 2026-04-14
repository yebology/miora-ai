"use client";

import { useState } from "react";
import { useAppKitAccount } from "@reown/appkit/react";
import type { SwapQuote } from "@/types/swap";
import { WalletGuardModal } from "@/components/ui/wallet-guard-modal";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { getTokensForChain, type Token } from "@/constants/tokens";
import {
  ArrowDownUp,
  Loader2,
  ChevronDown,
  Zap,
} from "lucide-react";

const CHAINS = [
  { value: "ethereum", label: "Ethereum" },
  { value: "arbitrum", label: "Arbitrum" },
  { value: "optimism", label: "Optimism" },
  { value: "base", label: "Base" },
  { value: "polygon", label: "Polygon" },
];

function TokenSelect({
  tokens,
  selected,
  onChange,
  exclude,
}: {
  tokens: Token[];
  selected: string;
  onChange: (addr: string) => void;
  exclude?: string;
}) {
  return (
    <div className="relative">
      <select
        value={selected}
        onChange={(e) => onChange(e.target.value)}
        className="h-10 w-full appearance-none rounded-lg border bg-card py-2 pl-3 pr-9 text-sm font-medium outline-none"
      >
        {tokens
          .filter((t) => t.address !== exclude)
          .map((t) => (
            <option key={t.address} value={t.address}>
              {t.symbol} — {t.name}
            </option>
          ))}
      </select>
      <ChevronDown className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
    </div>
  );
}

function shortenAddress(addr: string) {
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

export default function SwapPage() {
  const { isConnected } = useAppKitAccount();

  const [chain, setChain] = useState("ethereum");
  const [amount, setAmount] = useState("");
  const [loading, setLoading] = useState(false);
  const [quote, setQuote] = useState<SwapQuote | null>(null);
  const [showWalletModal, setShowWalletModal] = useState(false);

  const tokens = getTokensForChain(chain);
  const [inputMint, setInputMint] = useState(tokens[0]?.address || "");
  const [outputMint, setOutputMint] = useState(tokens[1]?.address || "");

  const handleChainChange = (newChain: string) => {
    setChain(newChain);
    const newTokens = getTokensForChain(newChain);
    setInputMint(newTokens[0]?.address || "");
    setOutputMint(newTokens[1]?.address || "");
    setQuote(null);
  };

  const handleSwapDirection = () => {
    setInputMint(outputMint);
    setOutputMint(inputMint);
    setQuote(null);
  };

  const inputToken = tokens.find((t) => t.address === inputMint);
  const outputToken = tokens.find((t) => t.address === outputMint);

  const handleGetQuote = async () => {
    if (!isConnected) {
      setShowWalletModal(true);
      return;
    }
    if (!amount || !inputMint || !outputMint) return;
    setLoading(true);
    setQuote(null);

    try {
      // TODO: POST /api/swap/quote
      await new Promise((r) => setTimeout(r, 1200));

      const inputDecimals = inputToken?.decimals || 18;
      const rawAmount = (
        parseFloat(amount) * Math.pow(10, inputDecimals)
      ).toString();
      const outputAmt = (parseFloat(amount) * 1847.5).toFixed(0);

      setQuote({
        chain,
        input_mint: inputMint,
        output_mint: outputMint,
        input_amount: rawAmount,
        output_amount: outputAmt,
        price_impact: "0.03",
        route: "Uniswap V3 → SushiSwap",
      });
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex flex-1 items-start justify-center px-6 py-10">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight">Swap Tokens</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Best price routing via 1inch on Base.
          </p>
        </div>

        <Card>
          <CardContent className="space-y-4 p-5">
            {/* Chain selector */}
            <div>
              <label className="mb-1.5 block text-xs text-muted-foreground">
                Chain
              </label>
              <div className="relative">
                <select
                  value={chain}
                  onChange={(e) => handleChainChange(e.target.value)}
                  className="h-10 w-full appearance-none rounded-lg border bg-card py-2 pl-3 pr-9 text-sm font-medium outline-none"
                >
                  {CHAINS.map((c) => (
                    <option key={c.value} value={c.value}>
                      {c.label}
                    </option>
                  ))}
                </select>
                <ChevronDown className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              </div>
            </div>

            {/* From */}
            <div>
              <label className="mb-1.5 block text-xs text-muted-foreground">
                From
              </label>
              <TokenSelect
                tokens={tokens}
                selected={inputMint}
                onChange={setInputMint}
                exclude={outputMint}
              />
              <Input
                type="number"
                placeholder="0.0"
                value={amount}
                onChange={(e) => {
                  setAmount(e.target.value);
                  setQuote(null);
                }}
                className="mt-2 h-12 font-mono text-lg"
              />
            </div>

            {/* Swap direction */}
            <div className="flex justify-center">
              <Button
                variant="ghost"
                size="icon"
                className="h-9 w-9 rounded-full border"
                onClick={handleSwapDirection}
              >
                <ArrowDownUp className="h-4 w-4" />
              </Button>
            </div>

            {/* To */}
            <div>
              <label className="mb-1.5 block text-xs text-muted-foreground">
                To
              </label>
              <TokenSelect
                tokens={tokens}
                selected={outputMint}
                onChange={setOutputMint}
                exclude={inputMint}
              />
              {quote && (
                <div className="mt-2 flex h-12 items-center rounded-lg border bg-muted/30 px-3 font-mono text-lg">
                  {quote.output_amount}
                  <span className="ml-2 text-sm text-muted-foreground">
                    {outputToken?.symbol}
                  </span>
                </div>
              )}
            </div>

            {/* Get Quote */}
            <Button
              className="w-full gap-2"
              onClick={handleGetQuote}
              disabled={loading || !amount || parseFloat(amount) <= 0}
            >
              {loading ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <Zap className="h-4 w-4" />
              )}
              {loading ? "Getting quote..." : "Get Quote"}
            </Button>

            {/* Quote details */}
            {quote && (
              <div className="space-y-2 rounded-lg border bg-muted/20 p-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Route</span>
                  <span className="font-medium">{quote.route}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Price Impact</span>
                  <span
                    className={
                      parseFloat(quote.price_impact || "0") > 1
                        ? "text-red-400"
                        : "text-green-400"
                    }
                  >
                    {quote.price_impact}%
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Provider</span>
                  <span className="font-medium">1inch</span>
                </div>
              </div>
            )}
          </CardContent>
        </Card>
        <WalletGuardModal open={showWalletModal} onOpenChange={setShowWalletModal} />
      </div>
    </div>
  );
}

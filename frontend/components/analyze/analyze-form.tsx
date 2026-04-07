"use client";

import { useState } from "react";
import { Search, Loader2, ChevronDown, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

const CHAINS = [
  { value: "ethereum", label: "Ethereum" },
  { value: "arbitrum", label: "Arbitrum" },
  { value: "optimism", label: "Optimism" },
  { value: "base", label: "Base" },
  { value: "polygon", label: "Polygon" },
  { value: "solana", label: "Solana" },
];

const EVM_LIMITS = [
  { value: 10, label: "10 txs", enabled: true },
  { value: 25, label: "25 txs", enabled: false },
  { value: 50, label: "50 txs", enabled: false },
  { value: 100, label: "100 txs", enabled: false },
];

const SVM_LIMITS = [
  { value: 20, label: "20 txs", enabled: true },
  { value: 50, label: "50 txs", enabled: false },
  { value: 100, label: "100 txs", enabled: false },
  { value: 200, label: "200 txs", enabled: false },
];

const EVM_CHAINS = ["ethereum", "arbitrum", "optimism", "base", "polygon"];

function getLimits(chain: string) {
  return EVM_CHAINS.includes(chain) ? EVM_LIMITS : SVM_LIMITS;
}

function getDefaultLimit(chain: string) {
  return EVM_CHAINS.includes(chain) ? 10 : 20;
}

type Props = {
  onAnalyze: (address: string, chain: string, limit: number) => void;
  loading: boolean;
  error?: string | null;
};

export function AnalyzeForm({ onAnalyze, loading, error }: Props) {
  const [address, setAddress] = useState("");
  const [chain, setChain] = useState("ethereum");
  const [limit, setLimit] = useState(10);

  const limits = getLimits(chain);

  const handleChainChange = (newChain: string) => {
    setChain(newChain);
    setLimit(getDefaultLimit(newChain));
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!address.trim()) return;
    onAnalyze(address.trim(), chain, limit);
  };

  return (
    <form onSubmit={handleSubmit} className="mx-auto max-w-2xl">
      <div className="flex flex-col gap-3 sm:flex-row">
        <Input
          placeholder="Paste wallet address (0x... or base58)"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          className="h-11 flex-1 font-mono text-sm"
        />
        <div className="relative">
          <select
            value={chain}
            onChange={(e) => handleChainChange(e.target.value)}
            className="h-11 w-full appearance-none rounded-lg border bg-card py-2 pl-3 pr-9 text-sm outline-none sm:w-auto"
          >
            {CHAINS.map((c) => (
              <option key={c.value} value={c.value}>
                {c.label}
              </option>
            ))}
          </select>
          <ChevronDown className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        </div>
        <Button type="submit" disabled={loading || !address.trim()} className="h-11 gap-2">
          {loading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Search className="h-4 w-4" />
          )}
          Analyze
        </Button>
      </div>

      {/* Limit selector */}
      <div className="mt-3 flex items-center justify-center gap-1.5">
        <span className="text-xs text-muted-foreground">Transactions:</span>
        <span className="group relative">
          <Info className="h-3.5 w-3.5 cursor-help text-muted-foreground/50 transition-colors hover:text-muted-foreground" />
          <span className="absolute bottom-full left-1/2 z-50 mb-2 hidden w-56 -translate-x-1/2 rounded-lg border bg-popover px-3 py-2 text-xs leading-relaxed text-popover-foreground shadow-lg group-hover:block">
            How many recent transactions to analyze. More transactions = more accurate scoring, but takes longer.
          </span>
        </span>
        {limits.map((l) => (
          <button
            key={l.value}
            type="button"
            disabled={!l.enabled}
            onClick={() => l.enabled && setLimit(l.value)}
            className={`rounded-md px-2.5 py-1 text-xs transition-colors ${
              limit === l.value
                ? "bg-primary text-primary-foreground"
                : l.enabled
                  ? "bg-muted text-muted-foreground hover:text-foreground"
                  : "cursor-not-allowed bg-muted/50 text-muted-foreground/30 line-through"
            }`}
          >
            {l.label}
            {!l.enabled && " 🔒"}
          </button>
        ))}
      </div>

      {error && (
        <p className="mt-3 text-center text-sm text-red-400">{error}</p>
      )}
    </form>
  );
}

"use client";

import { useState } from "react";
import { Search, Loader2, ChevronDown } from "lucide-react";
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

type Props = {
  onAnalyze: (address: string, chain: string) => void;
  loading: boolean;
  error?: string | null;
};

export function AnalyzeForm({ onAnalyze, loading, error }: Props) {
  const [address, setAddress] = useState("");
  const [chain, setChain] = useState("ethereum");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!address.trim()) return;
    onAnalyze(address.trim(), chain);
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
            onChange={(e) => setChain(e.target.value)}
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

      {error && (
        <p className="mt-3 text-center text-sm text-red-400">{error}</p>
      )}
    </form>
  );
}

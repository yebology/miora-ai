"use client";

import { useState } from "react";
import { Search, Loader2 } from "lucide-react";
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
};

export function AnalyzeForm({ onAnalyze, loading }: Props) {
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
        <select
          value={chain}
          onChange={(e) => setChain(e.target.value)}
          className="h-11 rounded-lg border bg-card px-3 text-sm outline-none"
        >
          {CHAINS.map((c) => (
            <option key={c.value} value={c.value}>
              {c.label}
            </option>
          ))}
        </select>
        <Button type="submit" disabled={loading || !address.trim()} className="h-11 gap-2">
          {loading ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Search className="h-4 w-4" />
          )}
          Analyze
        </Button>
      </div>
    </form>
  );
}

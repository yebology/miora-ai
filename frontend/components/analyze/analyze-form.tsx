"use client";

import { useState } from "react";
import { Search, Loader2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type Props = {
  onAnalyze: (address: string, chain: string) => void;
  loading: boolean;
  error?: string | null;
};

export function AnalyzeForm({ onAnalyze, loading, error }: Props) {
  const [address, setAddress] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!address.trim()) return;
    onAnalyze(address.trim(), "base");
  };

  return (
    <form onSubmit={handleSubmit} className="mx-auto max-w-2xl">
      <div className="flex flex-col gap-3 sm:flex-row">
        <Input
          placeholder="Paste wallet address (0x...)"
          value={address}
          onChange={(e) => setAddress(e.target.value)}
          className="h-11 flex-1 font-mono text-sm"
        />
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
        <div className="mt-3 flex items-center justify-center gap-2 rounded-lg bg-red-500/10 px-4 py-2.5 text-sm text-red-400">
          <Search className="h-4 w-4 shrink-0" />
          {error}
        </div>
      )}
    </form>
  );
}

"use client";

import { useState } from "react";
import type { WatchlistItem } from "@/types/watchlist";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Check, Bot, ChevronDown } from "lucide-react";
import { cn } from "@/lib/utils";

const CONDITIONS = [
  { id: "min_liquidity", label: "Min liquidity > $100k" },
  { id: "min_mcap", label: "Min market cap > $500k" },
  { id: "min_pair_age", label: "Min pair age > 6 hours" },
  { id: "min_volume", label: "Min 24h volume > $50k" },
];

function shortenAddress(addr: string) {
  if (addr.length <= 12) return addr;
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

type Props = {
  watchlist: WatchlistItem[];
  onCreate: (data: {
    bot_type: "wallet";
    target_wallet_address: string;
    target_wallet_chain: string;
    recommendation: string;
    budget: number;
    max_per_trade: number;
    conditions: string[];
  }) => void;
  creating: boolean;
  onCancel: () => void;
};

export function CreateBotForm({ watchlist, onCreate, creating, onCancel }: Props) {
  const [selectedWallet, setSelectedWallet] = useState<WatchlistItem | null>(null);
  const [budget, setBudget] = useState("");
  const [maxPerTrade, setMaxPerTrade] = useState("");
  const [selectedConditions, setSelectedConditions] = useState<Set<string>>(new Set());

  const isFullFollow = selectedWallet?.recommendation === "full_follow";

  const handleSelectWallet = (address: string) => {
    const item = watchlist.find((w) => w.wallet_address === address);
    if (item) {
      setSelectedWallet(item);
      setSelectedConditions(new Set(item.conditions || []));
    }
  };

  const toggleCondition = (id: string) => {
    setSelectedConditions((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedWallet || !budget || !maxPerTrade) return;
    onCreate({
      bot_type: "wallet",
      target_wallet_address: selectedWallet.wallet_address,
      target_wallet_chain: selectedWallet.chain,
      recommendation: selectedWallet.recommendation,
      budget: Number(budget),
      max_per_trade: Number(maxPerTrade),
      conditions: isFullFollow ? [] : Array.from(selectedConditions),
    });
  };

  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <Bot className="h-5 w-5 text-purple-400" />
          <h3 className="text-sm font-medium">Create Wallet Bot</h3>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <Label className="mb-1.5 block text-xs text-muted-foreground">Select wallet from watchlist</Label>
            <div className="relative">
              <select
                value={selectedWallet?.wallet_address || ""}
                onChange={(e) => handleSelectWallet(e.target.value)}
                className="h-10 w-full appearance-none rounded-lg border bg-card py-2 pl-3 pr-9 font-mono text-sm outline-none"
              >
                <option value="">Choose a wallet...</option>
                {watchlist.map((item) => (
                  <option key={item.wallet_address} value={item.wallet_address}>
                    {shortenAddress(item.wallet_address)} — {item.recommendation}
                  </option>
                ))}
              </select>
              <ChevronDown className="pointer-events-none absolute right-2.5 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            </div>
          </div>

          {selectedWallet && (
            <>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <Label className="text-xs text-muted-foreground">Budget (USD)</Label>
                  <Input type="number" min={0} step={10} value={budget}
                    onChange={(e) => setBudget(e.target.value)} placeholder="e.g. 500" className="mt-1" />
                </div>
                <div>
                  <Label className="text-xs text-muted-foreground">Max Per Trade (USD)</Label>
                  <Input type="number" min={0} step={5} value={maxPerTrade}
                    onChange={(e) => setMaxPerTrade(e.target.value)} placeholder="e.g. 50" className="mt-1" />
                </div>
              </div>

              {!isFullFollow && (
                <div>
                  <Label className="mb-2 block text-xs text-muted-foreground">Trade Conditions</Label>
                  <div className="space-y-2">
                    {CONDITIONS.map((c) => (
                      <button key={c.id} type="button" onClick={() => toggleCondition(c.id)}
                        className={cn("flex w-full items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition-all",
                          selectedConditions.has(c.id) ? "border-purple-500/40 bg-purple-500/10" : "border-border hover:border-muted-foreground/20")}>
                        <div className={cn("flex h-4 w-4 shrink-0 items-center justify-center rounded border",
                          selectedConditions.has(c.id) ? "border-purple-500 bg-purple-500 text-white" : "border-muted-foreground/30")}>
                          {selectedConditions.has(c.id) && <Check className="h-3 w-3" />}
                        </div>
                        <span>{c.label}</span>
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {isFullFollow && (
                <p className="rounded-lg bg-green-500/10 px-3 py-2 text-xs text-green-400">
                  Full Follow — bot will copy all trades without conditions.
                </p>
              )}

              <div className="flex gap-2">
                <Button type="button" variant="outline" className="flex-1" onClick={onCancel}>Cancel</Button>
                <Button type="submit" className="flex-1" disabled={creating || !budget || !maxPerTrade}>
                  {creating ? "Creating..." : "Create Bot"}
                </Button>
              </div>
            </>
          )}
        </form>
      </CardContent>
    </Card>
  );
}

"use client";

import { useState } from "react";
import type { AgentConfig } from "@/types/agent";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Check, Settings, Info } from "lucide-react";
import { cn } from "@/lib/utils";

const RISK_OPTIONS = [
  { value: "low", label: "Low", description: "Only high-liquidity, established tokens" },
  { value: "medium", label: "Medium", description: "Balanced — moderate liquidity tokens included" },
  { value: "high", label: "High", description: "Aggressive — includes newer, lower-liquidity tokens" },
] as const;

const CONDITIONS = [
  { id: "min_liquidity", label: "Min liquidity > $100k" },
  { id: "min_mcap", label: "Min market cap > $500k" },
  { id: "min_pair_age", label: "Min pair age > 6 hours" },
  { id: "min_volume", label: "Min 24h volume > $50k" },
];

type Props = {
  config: AgentConfig;
  onSave: (data: {
    budget: number;
    max_per_trade: number;
    risk_tolerance: string;
    min_score: number;
    conditions: string[];
  }) => void;
  saving: boolean;
};

export function AgentConfigForm({ config, onSave, saving }: Props) {
  const [budget, setBudget] = useState(config.budget || 0);
  const [maxPerTrade, setMaxPerTrade] = useState(config.max_per_trade || 0);
  const [riskTolerance, setRiskTolerance] = useState(config.risk_tolerance || "medium");
  const [minScore, setMinScore] = useState(config.min_score || 70);
  const [selectedConditions, setSelectedConditions] = useState<Set<string>>(
    new Set(config.conditions || [])
  );

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
    onSave({
      budget,
      max_per_trade: maxPerTrade,
      risk_tolerance: riskTolerance,
      min_score: minScore,
      conditions: Array.from(selectedConditions),
    });
  };

  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <Settings className="h-5 w-5 text-purple-400" />
          <h3 className="text-sm font-medium">Agent Configuration</h3>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          {/* Budget */}
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <Label className="text-xs text-muted-foreground">Total Budget (USD)</Label>
              <Input
                type="number"
                min={0}
                step={10}
                value={budget || ""}
                onChange={(e) => setBudget(Number(e.target.value))}
                placeholder="e.g. 500"
                className="mt-1"
              />
            </div>
            <div>
              <Label className="text-xs text-muted-foreground">Max Per Trade (USD)</Label>
              <Input
                type="number"
                min={0}
                step={5}
                value={maxPerTrade || ""}
                onChange={(e) => setMaxPerTrade(Number(e.target.value))}
                placeholder="e.g. 50"
                className="mt-1"
              />
            </div>
          </div>

          {/* Min Score */}
          <div>
            <div className="flex items-center gap-1.5">
              <Label className="text-xs text-muted-foreground">Minimum Wallet Score</Label>
              <span className="group relative">
                <Info className="h-3.5 w-3.5 cursor-help text-muted-foreground/50" />
                <span className="absolute bottom-full left-1/2 z-50 mb-2 hidden w-56 -translate-x-1/2 rounded-lg border bg-popover px-3 py-2 text-xs text-popover-foreground shadow-lg group-hover:block">
                  Agent will only copy trades from wallets with a Miora score at or above this threshold.
                </span>
              </span>
            </div>
            <Input
              type="number"
              min={0}
              max={100}
              value={minScore}
              onChange={(e) => setMinScore(Number(e.target.value))}
              className="mt-1 w-32"
            />
          </div>

          {/* Risk Tolerance */}
          <div>
            <Label className="mb-2 block text-xs text-muted-foreground">Risk Tolerance</Label>
            <div className="grid gap-2 sm:grid-cols-3">
              {RISK_OPTIONS.map((opt) => (
                <button
                  key={opt.value}
                  type="button"
                  onClick={() => setRiskTolerance(opt.value)}
                  className={cn(
                    "rounded-lg border px-3 py-2.5 text-left text-sm transition-all",
                    riskTolerance === opt.value
                      ? "border-purple-500/40 bg-purple-500/10"
                      : "border-border hover:border-muted-foreground/20"
                  )}
                >
                  <span className="font-medium">{opt.label}</span>
                  <span className="mt-0.5 block text-xs text-muted-foreground">{opt.description}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Conditions */}
          <div>
            <Label className="mb-2 block text-xs text-muted-foreground">Trade Conditions</Label>
            <div className="space-y-2">
              {CONDITIONS.map((c) => (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => toggleCondition(c.id)}
                  className={cn(
                    "flex w-full items-center gap-3 rounded-lg border px-4 py-2.5 text-left text-sm transition-all",
                    selectedConditions.has(c.id)
                      ? "border-purple-500/40 bg-purple-500/10"
                      : "border-border hover:border-muted-foreground/20"
                  )}
                >
                  <div
                    className={cn(
                      "flex h-4 w-4 shrink-0 items-center justify-center rounded border",
                      selectedConditions.has(c.id)
                        ? "border-purple-500 bg-purple-500 text-white"
                        : "border-muted-foreground/30"
                    )}
                  >
                    {selectedConditions.has(c.id) && <Check className="h-3 w-3" />}
                  </div>
                  <span>{c.label}</span>
                </button>
              ))}
            </div>
          </div>

          <Button type="submit" className="w-full" disabled={saving || budget <= 0 || maxPerTrade <= 0}>
            {saving ? "Saving..." : "Save Configuration"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

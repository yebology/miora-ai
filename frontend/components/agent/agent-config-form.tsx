"use client";

import { useState } from "react";
import type { AgentConfig } from "@/types/agent";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Check, Settings } from "lucide-react";
import { cn } from "@/lib/utils";

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
    conditions: string[];
    consensus_threshold?: number;
    consensus_window_min?: number;
    min_score?: number;
  }) => void;
  saving: boolean;
};

export function AgentConfigForm({ config, onSave, saving }: Props) {
  const [budget, setBudget] = useState(config.budget ? String(config.budget) : "");
  const [maxPerTrade, setMaxPerTrade] = useState(config.max_per_trade ? String(config.max_per_trade) : "");
  const [selectedConditions, setSelectedConditions] = useState<Set<string>>(new Set(config.conditions || []));
  const [threshold, setThreshold] = useState(config.consensus_threshold ? String(config.consensus_threshold) : "3");
  const [windowMin, setWindowMin] = useState(config.consensus_window_min || 60);
  const [windowUnit, setWindowUnit] = useState<"minutes" | "hours" | "days">(
    (config.consensus_window_min || 60) >= 1440 ? "days" : (config.consensus_window_min || 60) >= 60 ? "hours" : "minutes"
  );
  const [minScore, setMinScore] = useState(config.min_score ? String(config.min_score) : "70");

  const isConsensus = config.bot_type === "consensus";
  const isFullFollow = config.recommendation === "full_follow";

  const toggleCondition = (id: string) => {
    setSelectedConditions((prev) => { const next = new Set(prev); if (next.has(id)) next.delete(id); else next.add(id); return next; });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({
      budget: Number(budget),
      max_per_trade: Number(maxPerTrade),
      conditions: Array.from(selectedConditions),
      ...(isConsensus && {
        consensus_threshold: Number(threshold),
        consensus_window_min: windowMin,
        min_score: Number(minScore),
      }),
    });
  };

  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <Settings className="h-5 w-5 text-purple-400" />
          <h3 className="text-sm font-medium">Bot Configuration</h3>
        </div>

        <form onSubmit={handleSubmit} className="space-y-5">
          <div className="grid gap-4 sm:grid-cols-2">
            <div>
              <Label className="text-xs text-muted-foreground">Total Budget (USD)</Label>
              <Input type="number" min={0} step={10} value={budget}
                onChange={(e) => setBudget(e.target.value)} placeholder="e.g. 500" className="mt-1" />
            </div>
            <div>
              <Label className="text-xs text-muted-foreground">Max Per Trade (USD)</Label>
              <Input type="number" min={0} step={5} value={maxPerTrade}
                onChange={(e) => setMaxPerTrade(e.target.value)} placeholder="e.g. 50" className="mt-1" />
            </div>
          </div>

          {/* Wallet bot: conditions */}
          {!isConsensus && !isFullFollow && (
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

          {!isConsensus && isFullFollow && (
            <p className="rounded-lg bg-green-500/10 px-3 py-2 text-xs text-green-400">
              Full Follow — bot copies all trades without conditions.
            </p>
          )}

          {/* Consensus bot: threshold, window, min score */}
          {isConsensus && (
            <div className="space-y-4">
              <div>
                <Label className="text-xs text-muted-foreground">Min Wallet Score (0-100)</Label>
                <Input type="number" min={0} max={100} value={minScore}
                  onChange={(e) => setMinScore(e.target.value)} className="mt-1 w-32" />
              </div>
              <div className="grid gap-4 sm:grid-cols-2">
                <div>
                  <Label className="text-xs text-muted-foreground">Min wallets buying same token</Label>
                  <Input type="number" min={2} max={20} value={threshold}
                    onChange={(e) => setThreshold(e.target.value)} className="mt-1" />
                </div>
                <div>
                  <Label className="text-xs text-muted-foreground">Time window</Label>
                  <div className="mt-1 flex gap-2">
                    <Input type="number" min={1}
                      value={windowUnit === "minutes" ? windowMin : windowUnit === "hours" ? Math.round(windowMin / 60) : Math.round(windowMin / 1440)}
                      onChange={(e) => {
                        const val = Number(e.target.value);
                        if (windowUnit === "minutes") setWindowMin(val);
                        else if (windowUnit === "hours") setWindowMin(val * 60);
                        else setWindowMin(val * 1440);
                      }}
                      className="flex-1"
                    />
                    <select value={windowUnit} onChange={(e) => setWindowUnit(e.target.value as "minutes" | "hours" | "days")}
                      className="h-9 rounded-lg border bg-card px-2 text-sm">
                      <option value="minutes">Min</option>
                      <option value="hours">Hours</option>
                      <option value="days">Days</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          )}

          <Button type="submit" className="w-full" disabled={saving || !budget || !maxPerTrade}>
            {saving ? "Saving..." : "Save Configuration"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}

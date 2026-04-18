"use client";

import { useState } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Zap } from "lucide-react";

type Props = {
  onCreate: (data: {
    bot_type: "consensus";
    budget: number;
    max_per_trade: number;
    min_score: number;
    consensus_threshold: number;
    consensus_window_min: number;
  }) => void;
  creating: boolean;
  onCancel: () => void;
};

export function CreateConsensusForm({ onCreate, creating, onCancel }: Props) {
  const [budget, setBudget] = useState("");
  const [maxPerTrade, setMaxPerTrade] = useState("");
  const [minScore, setMinScore] = useState("70");
  const [threshold, setThreshold] = useState("3");
  const [windowMin, setWindowMin] = useState(60);
  const [windowUnit, setWindowUnit] = useState<"minutes" | "hours" | "days">("hours");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!budget || !maxPerTrade) return;
    onCreate({
      bot_type: "consensus",
      budget: Number(budget),
      max_per_trade: Number(maxPerTrade),
      min_score: Number(minScore),
      consensus_threshold: Number(threshold),
      consensus_window_min: windowMin,
    });
  };

  return (
    <Card className="border-yellow-500/20">
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <Zap className="h-5 w-5 text-yellow-400" />
          <h3 className="text-sm font-medium">Create Consensus Bot</h3>
          <span className="rounded-full bg-yellow-500/10 px-2 py-0.5 text-xs text-yellow-400">Premium</span>
        </div>

        <p className="mb-4 text-xs text-muted-foreground">
          Scans all wallets analyzed by Miora. Trades only when multiple high-score wallets buy the same token within a time window.
        </p>

        <form onSubmit={handleSubmit} className="space-y-4">
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

          <div className="flex gap-2">
            <Button type="button" variant="outline" className="flex-1" onClick={onCancel}>Cancel</Button>
            <Button type="submit" className="flex-1" disabled={creating || !budget || !maxPerTrade}>
              {creating ? "Creating..." : "Create Consensus Bot"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

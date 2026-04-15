"use client";

import type { AgentConfig } from "@/types/agent";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Bot, Play, Pause, Wallet, TrendingUp, DollarSign } from "lucide-react";
import { cn } from "@/lib/utils";

type Props = {
  config: AgentConfig;
  onStart: () => void;
  onPause: () => void;
  loading: boolean;
};

const STATUS_STYLES = {
  active: { color: "text-green-400", bg: "bg-green-500/10", label: "Active" },
  paused: { color: "text-yellow-400", bg: "bg-yellow-500/10", label: "Paused" },
  stopped: { color: "text-red-400", bg: "bg-red-500/10", label: "Stopped" },
};

export function AgentStatusCard({ config, onStart, onPause, loading }: Props) {
  const status = STATUS_STYLES[config.status] || STATUS_STYLES.paused;
  const remaining = config.budget - config.total_spent;

  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Bot className="h-5 w-5 text-purple-400" />
            <h3 className="text-sm font-medium">Agent Status</h3>
          </div>
          <span className={cn("rounded-full px-2.5 py-1 text-xs font-medium", status.bg, status.color)}>
            {status.label}
          </span>
        </div>

        <div className="mb-4 grid grid-cols-3 gap-3">
          <div className="rounded-lg bg-muted/50 p-3 text-center">
            <DollarSign className="mx-auto mb-1 h-4 w-4 text-muted-foreground" />
            <p className="text-lg font-semibold">${remaining.toFixed(0)}</p>
            <p className="text-xs text-muted-foreground">Remaining</p>
          </div>
          <div className="rounded-lg bg-muted/50 p-3 text-center">
            <TrendingUp className="mx-auto mb-1 h-4 w-4 text-muted-foreground" />
            <p className="text-lg font-semibold">{config.total_trades}</p>
            <p className="text-xs text-muted-foreground">Trades</p>
          </div>
          <div className="rounded-lg bg-muted/50 p-3 text-center">
            <Wallet className="mx-auto mb-1 h-4 w-4 text-muted-foreground" />
            <p className="text-lg font-semibold">${config.total_spent.toFixed(0)}</p>
            <p className="text-xs text-muted-foreground">Spent</p>
          </div>
        </div>

        {config.agent_wallet_address && (
          <div className="mb-4 rounded-lg bg-muted/30 px-3 py-2">
            <p className="text-xs text-muted-foreground">Agent Wallet</p>
            <p className="font-mono text-xs break-all">{config.agent_wallet_address}</p>
          </div>
        )}

        <div className="flex gap-2">
          {config.status === "active" ? (
            <Button variant="outline" className="flex-1 gap-1.5" onClick={onPause} disabled={loading}>
              <Pause className="h-3.5 w-3.5" />
              Pause Agent
            </Button>
          ) : (
            <Button className="flex-1 gap-1.5" onClick={onStart} disabled={loading || config.budget <= 0}>
              <Play className="h-3.5 w-3.5" />
              Start Agent
            </Button>
          )}
        </div>
      </CardContent>
    </Card>
  );
}

"use client";

import type { AgentTrade } from "@/types/agent";
import { Card, CardContent } from "@/components/ui/card";
import { History, ExternalLink, CheckCircle, XCircle, MinusCircle } from "lucide-react";
import { cn } from "@/lib/utils";

type Props = {
  trades: AgentTrade[];
};

const STATUS_ICON = {
  executed: <CheckCircle className="h-4 w-4 text-green-400" />,
  failed: <XCircle className="h-4 w-4 text-red-400" />,
  skipped: <MinusCircle className="h-4 w-4 text-yellow-400" />,
};

export function AgentTradeHistory({ trades }: Props) {
  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <History className="h-5 w-5 text-purple-400" />
          <h3 className="text-sm font-medium">Trade History</h3>
        </div>

        {trades.length === 0 ? (
          <div className="py-8 text-center">
            <History className="mx-auto mb-2 h-8 w-8 text-muted-foreground/30" />
            <p className="text-sm text-muted-foreground">No trades yet. Start the agent to begin.</p>
          </div>
        ) : (
          <div className="space-y-2">
            {trades.map((trade) => (
              <div
                key={trade.id}
                className={cn(
                  "flex items-start gap-3 rounded-lg border px-4 py-3",
                  trade.status === "executed" && "border-green-500/20",
                  trade.status === "failed" && "border-red-500/20",
                  trade.status === "skipped" && "border-yellow-500/20"
                )}
              >
                {STATUS_ICON[trade.status]}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <span className="text-sm font-medium">{trade.token_symbol}</span>
                    <span className="text-xs text-muted-foreground">${trade.amount_usd.toFixed(2)}</span>
                    <span className={cn(
                      "rounded-full px-1.5 py-0.5 text-xs",
                      trade.status === "executed" && "bg-green-500/10 text-green-400",
                      trade.status === "failed" && "bg-red-500/10 text-red-400",
                      trade.status === "skipped" && "bg-yellow-500/10 text-yellow-400"
                    )}>
                      {trade.status}
                    </span>
                  </div>
                  <p className="mt-0.5 text-xs text-muted-foreground truncate">{trade.reason}</p>
                  <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                    <span>Source: {trade.source_wallet.slice(0, 8)}... (score {trade.source_score})</span>
                    {trade.tx_hash && (
                      <a
                        href={`https://sepolia.basescan.org/tx/${trade.tx_hash}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex items-center gap-0.5 text-purple-400 hover:underline"
                      >
                        <ExternalLink className="h-3 w-3" />
                        tx
                      </a>
                    )}
                  </div>
                </div>
                <span className="text-xs text-muted-foreground whitespace-nowrap">
                  {new Date(trade.created_at).toLocaleDateString()}
                </span>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}

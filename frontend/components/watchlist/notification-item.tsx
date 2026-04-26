"use client";

import { useState } from "react";
import type { Notification } from "@/types/watchlist";
import { cn } from "@/lib/utils";
import { ArrowDownLeft, ArrowUpRight, ChevronDown } from "lucide-react";

type Props = {
  notification: Notification;
};

function shortenAddress(addr: string) {
  if (addr.length <= 12) return addr;
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

function formatMoney(n: number): string {
  if (n >= 1_000_000_000) return `${(n / 1_000_000_000).toFixed(1)}B`;
  if (n >= 1_000_000) return `${(n / 1_000_000).toFixed(1)}M`;
  if (n >= 1_000) return `${(n / 1_000).toFixed(0)}k`;
  return `${n.toFixed(0)}`;
}

function timeAgo(dateStr: string): string {
  const diff = Date.now() - new Date(dateStr).getTime();
  const mins = Math.floor(diff / 60000);
  if (mins < 60) return `${mins}m ago`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  return `${days}d ago`;
}

export function NotificationItem({ notification: n }: Props) {
  const isBuy = n.direction === "in";
  const [expanded, setExpanded] = useState(false);

  return (
    <div
      className={cn(
        "rounded-lg border transition-colors",
        n.read ? "bg-card" : "border-primary/20 bg-primary/5"
      )}
    >
      <div className="flex items-start gap-3 px-4 py-3">
        <div
          className={cn(
            "mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-full",
            isBuy ? "bg-green-500/10 text-green-400" : "bg-red-500/10 text-red-400"
          )}
        >
          {isBuy ? (
            <ArrowDownLeft className="h-4 w-4" />
          ) : (
            <ArrowUpRight className="h-4 w-4" />
          )}
        </div>

        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="text-sm font-medium">
              {shortenAddress(n.wallet_address)}{" "}
              <span className={isBuy ? "text-green-400" : "text-red-400"}>
                {isBuy ? "bought" : "sold"}
              </span>{" "}
              <span className="font-mono">{n.token_symbol}</span>
            </span>
          </div>

          <div className="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
            <span>Amount: {n.value}</span>
            <span>Liq: {formatMoney(n.liquidity)}</span>
            <span>MCap: {formatMoney(n.market_cap)}</span>
            <span className="capitalize">{n.chain}</span>
          </div>
        </div>

        <div className="flex shrink-0 flex-col items-end gap-1.5">
          <span className="text-xs text-muted-foreground">
            {timeAgo(n.created_at)}
          </span>
          {n.ai_assessment && (
            <button
              onClick={() => setExpanded(!expanded)}
              className="flex items-center gap-1 rounded-md px-1.5 py-0.5 text-xs text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              AI Insight
              <ChevronDown className={cn("h-3 w-3 transition-transform", expanded && "rotate-180")} />
            </button>
          )}
        </div>
      </div>

      {expanded && n.ai_assessment && (
        <div className="border-t px-4 py-2.5">
          <p className="text-xs leading-relaxed text-muted-foreground">
            {n.ai_assessment}
          </p>
        </div>
      )}
    </div>
  );
}

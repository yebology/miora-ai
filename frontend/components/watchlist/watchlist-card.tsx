"use client";

import type { WatchlistItem } from "@/types/watchlist";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Trash2, ExternalLink, Bell, BellOff } from "lucide-react";
import { cn } from "@/lib/utils";
import Link from "next/link";

type Props = {
  item: WatchlistItem;
  onUnfollow: (address: string) => void;
  onToggleNotify: (address: string) => void;
};

const CONDITION_LABELS: Record<string, string> = {
  min_liquidity: "Liquidity > $100k",
  min_mcap: "MCap > $500k",
  min_pair_age: "Pair age > 6h",
  min_volume: "Volume > $50k",
};

const REC_CONFIG = {
  full_follow: { label: "Full Follow", className: "text-green-400" },
  conditional_follow: { label: "Conditional", className: "text-yellow-400" },
  avoid: { label: "Avoid", className: "text-red-400" },
} as const;

function shortenAddress(addr: string) {
  if (addr.length <= 12) return addr;
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

export function WatchlistCard({ item, onUnfollow, onToggleNotify }: Props) {
  const rec = REC_CONFIG[item.recommendation as keyof typeof REC_CONFIG] || REC_CONFIG.conditional_follow;

  return (
    <Card>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-3">
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2">
              <span className="font-mono text-sm font-medium">
                {shortenAddress(item.wallet_address)}
              </span>
              <span className="rounded-full bg-muted px-2 py-0.5 text-xs capitalize">
                {item.chain}
              </span>
              <span className={cn("text-xs font-medium", rec.className)}>
                {rec.label}
              </span>
            </div>

            {item.conditions.length > 0 && (
              <div className="mt-2 flex flex-wrap gap-1.5">
                {item.conditions.map((c) => (
                  <span
                    key={c}
                    className="rounded-full border border-yellow-500/20 bg-yellow-500/5 px-2 py-0.5 text-xs text-yellow-400"
                  >
                    {CONDITION_LABELS[c] || c}
                  </span>
                ))}
              </div>
            )}

            <div className="mt-2 text-xs text-muted-foreground">
              Following since {new Date(item.created_at).toLocaleDateString()}
            </div>
          </div>

          <div className="flex shrink-0 gap-1">
            <Button
              variant="ghost"
              size="icon"
              className={cn(
                "h-8 w-8",
                item.email_notify
                  ? "text-primary hover:text-primary/80"
                  : "text-muted-foreground hover:text-foreground"
              )}
              onClick={() => onToggleNotify(item.wallet_address)}
              title={item.email_notify ? "Notifications on" : "Notifications off"}
            >
              {item.email_notify ? (
                <Bell className="h-3.5 w-3.5" />
              ) : (
                <BellOff className="h-3.5 w-3.5" />
              )}
            </Button>
            <Link href={`/watchlist/${item.chain}/${item.wallet_address}`}>
              <Button variant="ghost" size="icon" className="h-8 w-8">
                <ExternalLink className="h-3.5 w-3.5" />
              </Button>
            </Link>
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 text-red-400 hover:text-red-300"
              onClick={() => onUnfollow(item.wallet_address)}
            >
              <Trash2 className="h-3.5 w-3.5" />
            </Button>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

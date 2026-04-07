"use client";

import { useState } from "react";
import type { Condition } from "@/types/wallet";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Eye, Check, ShieldAlert } from "lucide-react";
import { cn } from "@/lib/utils";

type Props = {
  conditions: Condition[];
  address: string;
  chain: string;
};

export function ConditionsCard({ conditions, address, chain }: Props) {
  const [selected, setSelected] = useState<Set<string>>(
    new Set(conditions.map((c) => c.id))
  );
  const [followed, setFollowed] = useState(false);

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleFollow = async () => {
    // TODO: Replace with real API call
    // await fetch(`${API_URL}/api/watchlist/follow`, {
    //   method: "POST",
    //   headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
    //   body: JSON.stringify({
    //     wallet_address: address,
    //     chain,
    //     conditions: Array.from(selected),
    //   }),
    // });

    setFollowed(true);
  };

  return (
    <Card className="border-yellow-500/20">
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <ShieldAlert className="h-5 w-5 text-yellow-400" />
          <h3 className="text-sm font-medium">Follow with Conditions</h3>
        </div>

        <p className="mb-4 text-xs text-muted-foreground">
          This wallet is conditionally recommended. Select which conditions to
          apply — you&apos;ll only get notified when trades match your selected
          filters.
        </p>

        <div className="mb-4 space-y-2">
          {conditions.map((c) => (
            <button
              key={c.id}
              onClick={() => toggle(c.id)}
              disabled={followed}
              className={cn(
                "flex w-full items-center gap-3 rounded-lg border px-4 py-3 text-left text-sm transition-all",
                selected.has(c.id)
                  ? "border-yellow-500/40 bg-yellow-500/10"
                  : "border-border bg-card hover:border-muted-foreground/20"
              )}
            >
              <div
                className={cn(
                  "flex h-5 w-5 shrink-0 items-center justify-center rounded border transition-colors",
                  selected.has(c.id)
                    ? "border-yellow-500 bg-yellow-500 text-black"
                    : "border-muted-foreground/30"
                )}
              >
                {selected.has(c.id) && <Check className="h-3 w-3" />}
              </div>
              <span
                className={cn(
                  "text-sm",
                  selected.has(c.id)
                    ? "text-foreground"
                    : "text-muted-foreground"
                )}
              >
                {c.label}
                <span className="mt-0.5 block text-xs text-muted-foreground">
                  {c.description}
                </span>
              </span>
            </button>
          ))}
        </div>

        {followed ? (
          <div className="flex items-center gap-2 rounded-lg bg-green-500/10 px-4 py-3 text-sm text-green-400">
            <Check className="h-4 w-4" />
            Following with {selected.size} condition
            {selected.size !== 1 ? "s" : ""}
          </div>
        ) : (
          <Button
            onClick={handleFollow}
            disabled={selected.size === 0}
            className="w-full gap-2"
          >
            <Eye className="h-4 w-4" />
            Follow Wallet
            {selected.size > 0 && (
              <span className="rounded-full bg-primary-foreground/20 px-1.5 py-0.5 text-xs">
                {selected.size}
              </span>
            )}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}

"use client";

import { useState } from "react";
import type { Condition } from "@/types/wallet";
import { useAuth } from "@/components/providers/auth-provider";
import { AuthGuardModal } from "@/components/ui/auth-guard-modal";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Eye, Check, ShieldAlert, Info } from "lucide-react";
import { cn } from "@/lib/utils";
import { followWallet } from "@/api/watchlist/connector";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription,
} from "@/components/ui/dialog";

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
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const { user } = useAuth();

  const toggle = (id: string) => {
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) next.delete(id);
      else next.add(id);
      return next;
    });
  };

  const handleFollow = () => {
    if (!user) {
      setShowAuthModal(true);
      return;
    }
    setShowConfirm(true);
  };

  const confirmFollow = async () => {
    setFollowLoading(true);
    try {
      await followWallet(user!.walletAddress, {
        wallet_address: address,
        chain,
        recommendation: "conditional_follow",
        conditions: Array.from(selected),
        email_notify: false,
      });
      setFollowed(true);
    } catch {
      // Silently fail
    } finally {
      setFollowLoading(false);
      setShowConfirm(false);
    }
  };

  return (
    <Card className="border-yellow-500/20">
      <CardContent className="p-5">
        <div className="mb-4 flex items-center gap-2">
          <ShieldAlert className="h-5 w-5 text-yellow-400" />
          <h3 className="text-sm font-medium">Follow with Conditions</h3>
          <span className="group relative">
            <Info className="h-4 w-4 cursor-help text-muted-foreground/50 transition-colors hover:text-muted-foreground" />
            <span className="absolute top-full left-0 z-50 mt-2 hidden w-72 rounded-lg border bg-popover px-3 py-2 text-xs leading-relaxed text-popover-foreground shadow-lg group-hover:block">
              Follow this wallet to get notified when it trades — but only for tokens that match your selected conditions. This filters out risky trades so you only see the ones worth acting on.
            </span>
          </span>
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

        {/* Follow Confirmation */}
        <Dialog open={showConfirm} onOpenChange={setShowConfirm}>
          <DialogContent className="max-w-sm">
            <DialogHeader className="items-center text-center">
              <ShieldAlert className="mb-2 h-10 w-10 text-yellow-400" />
              <DialogTitle>Follow with {selected.size} condition{selected.size !== 1 ? "s" : ""}?</DialogTitle>
              <DialogDescription>
                You&apos;ll only get notified when this wallet trades tokens that match your selected conditions.
              </DialogDescription>
            </DialogHeader>
            <div className="flex gap-2">
              <Button variant="outline" className="flex-1" onClick={() => setShowConfirm(false)}>Cancel</Button>
              <Button className="flex-1" onClick={confirmFollow} disabled={followLoading}>
                {followLoading ? "Following..." : "Follow"}
              </Button>
            </div>
          </DialogContent>
        </Dialog>

        <AuthGuardModal open={showAuthModal} onOpenChange={setShowAuthModal} />
      </CardContent>
    </Card>
  );
}

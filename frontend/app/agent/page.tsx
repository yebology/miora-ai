"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/components/providers/auth-provider";
import { AuthGuardModal } from "@/components/ui/auth-guard-modal";
import { CreateBotForm } from "@/components/agent/create-bot-form";
import { CreateConsensusForm } from "@/components/agent/create-consensus-form";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Bot, Play, Pause, Trash2, ExternalLink, AlertTriangle, Zap, Wallet, Loader2 } from "lucide-react";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";
import Link from "next/link";
import { listBots, createBot, deleteBot, startBot, pauseBot } from "@/api/agent/connector";
import { getWatchlist } from "@/api/watchlist/connector";
import type { AgentConfig } from "@/api/agent/validation";
import type { WatchlistItem } from "@/api/watchlist/validation";

function shortenAddress(addr: string) {
  if (addr.length <= 12) return addr;
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

const STATUS_STYLES = {
  active: { color: "text-green-400", bg: "bg-green-500/10", label: "Active" },
  paused: { color: "text-yellow-400", bg: "bg-yellow-500/10", label: "Paused" },
  stopped: { color: "text-red-400", bg: "bg-red-500/10", label: "Stopped" },
};

export default function AgentPage() {
  const { user, isConnected } = useAuth();
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [bots, setBots] = useState<AgentConfig[]>([]);
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>([]);
  const [createMode, setCreateMode] = useState<"none" | "wallet" | "consensus">("none");
  const [creating, setCreating] = useState(false);
  const [deleteTarget, setDeleteTarget] = useState<number | null>(null);
  const [toggleTarget, setToggleTarget] = useState<AgentConfig | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!isConnected || !user) return;
    setLoading(true);
    Promise.all([
      listBots(user.walletAddress).catch(() => []),
      getWatchlist(user.walletAddress).catch(() => []),
    ])
      .then(([botsData, watchlistData]) => {
        setBots(botsData);
        setWatchlist(watchlistData);
      })
      .finally(() => setLoading(false));
  }, [isConnected, user]);

  const handleCreateWallet = async (data: any) => {
    if (!user) return;
    setCreating(true);
    try {
      const bot = await createBot(user.walletAddress, {
        ...data,
        bot_type: "wallet",
      });
      setBots((prev) => [bot, ...prev]);
      setCreateMode("none");
    } catch {
      // Silently fail
    } finally {
      setCreating(false);
    }
  };

  const handleCreateConsensus = async (data: any) => {
    if (!user) return;
    setCreating(true);
    try {
      const bot = await createBot(user.walletAddress, {
        ...data,
        bot_type: "consensus",
        conditions: [],
      });
      setBots((prev) => [bot, ...prev]);
      setCreateMode("none");
    } catch {
      // Silently fail
    } finally {
      setCreating(false);
    }
  };

  const handleDelete = async (id: number) => {
    if (!user) return;
    try {
      await deleteBot(user.walletAddress, id);
      setBots((prev) => prev.filter((b) => b.id !== id));
    } catch {
      // Silently fail
    }
    setDeleteTarget(null);
  };

  const handleToggle = async (id: number) => {
    if (!user) return;
    const bot = bots.find((b) => b.id === id);
    if (!bot) return;
    try {
      const updated = bot.status === "active"
        ? await pauseBot(user.walletAddress, id)
        : await startBot(user.walletAddress, id);
      setBots((prev) => prev.map((b) => (b.id === id ? updated : b)));
    } catch {
      // Silently fail
    }
    setToggleTarget(null);
  };

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight">AI Trading Bots</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Create bots that auto-trade based on wallets you&apos;ve analyzed.
          </p>
        </div>

        {!isConnected ? (
          <>
            <div className="py-16 text-center">
              <Bot className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
              <p className="mb-4 text-sm text-muted-foreground">Connect your wallet to create and manage trading bots.</p>
              <button onClick={() => setShowAuthModal(true)} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90">Connect Wallet</button>
            </div>
            <AuthGuardModal open={showAuthModal} onOpenChange={setShowAuthModal} />
          </>
        ) : loading ? (
          <div className="flex justify-center py-16">
            <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
          </div>
        ) : (
          <>
            {/* Create forms */}
            {createMode === "wallet" && (
              <div className="mb-6">
                <CreateBotForm watchlist={watchlist} onCreate={handleCreateWallet} creating={creating} onCancel={() => setCreateMode("none")} />
              </div>
            )}
            {createMode === "consensus" && (
              <div className="mb-6">
                <CreateConsensusForm onCreate={handleCreateConsensus} creating={creating} onCancel={() => setCreateMode("none")} />
              </div>
            )}

            {/* Create buttons */}
            {createMode === "none" && (
              <div className="mb-6 grid gap-3 sm:grid-cols-2">
                <Card className="cursor-pointer border-dashed transition-colors hover:border-purple-500/30" onClick={() => setCreateMode("wallet")}>
                  <CardContent className="flex items-center gap-3 p-4">
                    <Wallet className="h-5 w-5 text-purple-400" />
                    <div>
                      <p className="text-sm font-medium">Wallet Bot</p>
                      <p className="text-xs text-muted-foreground">Copy trades from a specific wallet</p>
                    </div>
                  </CardContent>
                </Card>
                <Card className="cursor-pointer border-dashed border-yellow-500/20 transition-colors hover:border-yellow-500/40" onClick={() => setCreateMode("consensus")}>
                  <CardContent className="flex items-center gap-3 p-4">
                    <Zap className="h-5 w-5 text-yellow-400" />
                    <div>
                      <p className="text-sm font-medium">Consensus Bot <span className="ml-1 rounded-full bg-yellow-500/10 px-1.5 py-0.5 text-xs text-yellow-400">Premium</span></p>
                      <p className="text-xs text-muted-foreground">Trade when multiple wallets buy same token</p>
                    </div>
                  </CardContent>
                </Card>
              </div>
            )}

            {/* Bot list */}
            {bots.length === 0 ? (
              <div className="py-16 text-center">
                <Bot className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
                <p className="text-sm text-muted-foreground">No bots yet. Create one above.</p>
              </div>
            ) : (
              <div className="space-y-3">
                {bots.map((bot) => {
                  const status = STATUS_STYLES[bot.status] || STATUS_STYLES.paused;
                  const remaining = bot.budget - bot.total_spent;
                  const isConsensus = bot.bot_type === "consensus";

                  return (
                    <Card key={bot.id} className={isConsensus ? "border-yellow-500/10" : ""}>
                      <CardContent className="p-4">
                        <div className="flex items-start justify-between gap-3">
                          <div className="min-w-0 flex-1">
                            <div className="flex items-center gap-2">
                              {isConsensus ? (
                                <Zap className="h-4 w-4 text-yellow-400" />
                              ) : (
                                <Wallet className="h-4 w-4 text-purple-400" />
                              )}
                              <span className="text-sm font-medium">
                                {isConsensus ? "Consensus Bot" : shortenAddress(bot.target_wallet_address || "")}
                              </span>
                              {!isConsensus && bot.target_wallet_score && (
                                <span className="rounded-full bg-muted px-2 py-0.5 text-xs">Score {bot.target_wallet_score}</span>
                              )}
                              {isConsensus && (
                                <span className="rounded-full bg-yellow-500/10 px-2 py-0.5 text-xs text-yellow-400">
                                  {bot.consensus_threshold}+ wallets · {bot.min_score}+ score
                                </span>
                              )}
                              <span className={cn("rounded-full px-2 py-0.5 text-xs", status.bg, status.color)}>{status.label}</span>
                            </div>
                            <div className="mt-2 flex gap-4 text-xs text-muted-foreground">
                              <span>${remaining.toFixed(0)} remaining</span>
                              <span>{bot.total_trades} trades</span>
                              <span>${bot.total_spent.toFixed(0)} spent</span>
                            </div>
                          </div>
                          <div className="flex shrink-0 gap-1">
                            <Button variant="ghost" size="icon" className="h-8 w-8" onClick={() => setToggleTarget(bot)}>
                              {bot.status === "active" ? <Pause className="h-3.5 w-3.5 text-yellow-400" /> : <Play className="h-3.5 w-3.5 text-green-400" />}
                            </Button>
                            <Link href={`/agent/${bot.id}`}>
                              <Button variant="ghost" size="icon" className="h-8 w-8"><ExternalLink className="h-3.5 w-3.5" /></Button>
                            </Link>
                            <Button variant="ghost" size="icon" className="h-8 w-8 text-red-400 hover:text-red-300" onClick={() => setDeleteTarget(bot.id)}>
                              <Trash2 className="h-3.5 w-3.5" />
                            </Button>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            )}

            {/* Toggle confirmation */}
            <Dialog open={toggleTarget !== null} onOpenChange={() => setToggleTarget(null)}>
              <DialogContent className="max-w-sm">
                <DialogHeader className="items-center text-center">
                  {toggleTarget?.status === "active" ? <Pause className="mb-2 h-10 w-10 text-yellow-400" /> : <Play className="mb-2 h-10 w-10 text-green-400" />}
                  <DialogTitle>{toggleTarget?.status === "active" ? "Pause this bot?" : "Start this bot?"}</DialogTitle>
                  <DialogDescription>
                    {toggleTarget?.status === "active" ? "The bot will stop monitoring and trading." : "The bot will start auto-trading when conditions are met."}
                  </DialogDescription>
                </DialogHeader>
                <div className="flex gap-2">
                  <Button variant="outline" className="flex-1" onClick={() => setToggleTarget(null)}>Cancel</Button>
                  <Button className="flex-1" onClick={() => toggleTarget && handleToggle(toggleTarget.id)}>
                    {toggleTarget?.status === "active" ? "Pause" : "Start"}
                  </Button>
                </div>
              </DialogContent>
            </Dialog>

            {/* Delete confirmation */}
            <Dialog open={deleteTarget !== null} onOpenChange={() => setDeleteTarget(null)}>
              <DialogContent className="max-w-sm">
                <DialogHeader className="items-center text-center">
                  <AlertTriangle className="mb-2 h-10 w-10 text-red-400" />
                  <DialogTitle>Delete this bot?</DialogTitle>
                  <DialogDescription>This will permanently remove the bot and its trade history.</DialogDescription>
                </DialogHeader>
                <div className="flex gap-2">
                  <Button variant="outline" className="flex-1" onClick={() => setDeleteTarget(null)}>Cancel</Button>
                  <Button variant="destructive" className="flex-1" onClick={() => deleteTarget && handleDelete(deleteTarget)}>Delete</Button>
                </div>
              </DialogContent>
            </Dialog>
          </>
        )}
      </div>
    </div>
  );
}

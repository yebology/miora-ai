"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import { useAuth } from "@/components/providers/auth-provider";
import { AgentStatusCard } from "@/components/agent/agent-status-card";
import { AgentConfigForm } from "@/components/agent/agent-config-form";
import { AgentTradeHistory } from "@/components/agent/agent-trade-history";
import { ArrowLeft, Loader2 } from "lucide-react";
import Link from "next/link";
import { getBot, updateBot, startBot, pauseBot, getBotTrades } from "@/api/agent/connector";
import type { AgentConfig, AgentTrade } from "@/api/agent/validation";

export default function BotDetailPage() {
  const params = useParams();
  const id = Number(params.id);
  const { user, isConnected } = useAuth();

  const [config, setConfig] = useState<AgentConfig | null>(null);
  const [trades, setTrades] = useState<AgentTrade[]>([]);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    if (!isConnected || !user) return;
    setLoading(true);
    Promise.all([
      getBot(user.walletAddress, id),
      getBotTrades(user.walletAddress, id, 50),
    ])
      .then(([botData, tradesData]) => {
        setConfig(botData);
        setTrades(tradesData);
      })
      .catch(() => setConfig(null))
      .finally(() => setLoading(false));
  }, [isConnected, user, id]);

  const handleSaveConfig = async (data: {
    budget: number;
    max_per_trade: number;
    conditions: string[];
    consensus_threshold?: number;
    consensus_window_min?: number;
    min_score?: number;
  }) => {
    if (!user || !config) return;
    setSaving(true);
    try {
      const updated = await updateBot(user.walletAddress, config.id, data);
      setConfig(updated);
    } catch {
      // Silently fail
    } finally {
      setSaving(false);
    }
  };

  const handleStart = async () => {
    if (!user || !config) return;
    setActionLoading(true);
    try {
      const updated = await startBot(user.walletAddress, config.id);
      setConfig(updated);
    } catch {
      // Silently fail
    } finally {
      setActionLoading(false);
    }
  };

  const handlePause = async () => {
    if (!user || !config) return;
    setActionLoading(true);
    try {
      const updated = await pauseBot(user.walletAddress, config.id);
      setConfig(updated);
    } catch {
      // Silently fail
    } finally {
      setActionLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (!config) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <div className="text-center">
          <p className="text-sm text-muted-foreground">Bot not found.</p>
          <Link href="/agent" className="mt-4 inline-block text-sm text-primary hover:underline">
            Back to Bots
          </Link>
        </div>
      </div>
    );
  }

  const isConsensus = config.bot_type === "consensus";

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <div className="mb-6">
          <Link href="/agent" className="flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground">
            <ArrowLeft className="h-4 w-4" />
            Back to Bots
          </Link>
        </div>

        <div className="mb-6">
          <h1 className="text-xl font-bold tracking-tight">
            {isConsensus ? "Consensus Bot" : `Bot — ${(config.target_wallet_address || "").slice(0, 8)}...`}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {isConsensus
              ? `${config.consensus_threshold}+ wallets · score ${config.min_score}+`
              : `Score ${config.target_wallet_score} · ${config.recommendation} · ${config.target_wallet_chain}`
            }
          </p>
        </div>

        <div className="space-y-6">
          <AgentStatusCard config={config} onStart={handleStart} onPause={handlePause} loading={actionLoading} />
          <AgentConfigForm config={config} onSave={handleSaveConfig} saving={saving} />
          <AgentTradeHistory trades={trades} />
        </div>
      </div>
    </div>
  );
}

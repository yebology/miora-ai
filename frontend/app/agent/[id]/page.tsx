"use client";

import { useState } from "react";
import { useParams } from "next/navigation";
import type { AgentConfig, AgentTrade } from "@/types/agent";
import { useAuth } from "@/components/providers/auth-provider";
import { AgentStatusCard } from "@/components/agent/agent-status-card";
import { AgentConfigForm } from "@/components/agent/agent-config-form";
import { AgentTradeHistory } from "@/components/agent/agent-trade-history";
import { ArrowLeft } from "lucide-react";
import Link from "next/link";

// Dummy bots — pick based on URL id
const DUMMY_BOTS: Record<string, AgentConfig> = {
  "1": {
    id: 1, user_id: 1, bot_type: "wallet",
    target_wallet_address: "0x1234567890abcdef1234567890abcdef12345678",
    target_wallet_chain: "base", target_wallet_score: 87,
    recommendation: "conditional_follow",
    budget: 500, max_per_trade: 50, conditions: ["min_liquidity", "min_mcap"],
    status: "active", agent_wallet_address: "0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18",
    total_spent: 120, total_trades: 4,
    created_at: "2025-12-01T10:00:00Z", updated_at: "2025-12-10T14:00:00Z",
  },
  "2": {
    id: 2, user_id: 1, bot_type: "consensus",
    budget: 300, max_per_trade: 30, conditions: [],
    status: "paused", agent_wallet_address: "", total_spent: 0, total_trades: 0,
    consensus_threshold: 3, consensus_window_min: 60, min_score: 75,
    created_at: "2025-12-05T14:30:00Z", updated_at: "2025-12-05T14:30:00Z",
  },
  "3": {
    id: 3, user_id: 1, bot_type: "wallet",
    target_wallet_address: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
    target_wallet_chain: "base", target_wallet_score: 92,
    recommendation: "full_follow",
    budget: 200, max_per_trade: 20, conditions: [],
    status: "paused", agent_wallet_address: "", total_spent: 0, total_trades: 0,
    created_at: "2025-12-08T09:00:00Z", updated_at: "2025-12-08T09:00:00Z",
  },
};

const DUMMY_TRADES: AgentTrade[] = [
  {
    id: 1, agent_config_id: 1,
    source_wallet: "0x1234567890abcdef1234567890abcdef12345678",
    source_score: 87, token_address: "0x6982508145454ce325ddbe47a25d4ec3d2311933",
    token_symbol: "PEPE", direction: "buy", amount_usd: 45,
    tx_hash: "0xabc123def456789012345678901234567890abcdef1234567890abcdef123456",
    status: "executed",
    reason: "Bought PEPE because wallet 0x1234...5678 (score 87) bought it.",
    risk_assessment: "Low risk.", created_at: "2025-12-10T15:30:00Z",
  },
  {
    id: 2, agent_config_id: 1,
    source_wallet: "0x1234567890abcdef1234567890abcdef12345678",
    source_score: 87, token_address: "0x0000000000000000000000000000000000001337",
    token_symbol: "NEWTOKEN", direction: "buy", amount_usd: 50, tx_hash: "",
    status: "skipped",
    reason: "Token liquidity $8k — below min_liquidity condition.",
    risk_assessment: "High risk.", created_at: "2025-12-09T11:00:00Z",
  },
  {
    id: 3, agent_config_id: 1,
    source_wallet: "0x1234567890abcdef1234567890abcdef12345678",
    source_score: 87, token_address: "0x514910771af9ca656af840dff83e8264ecf986ca",
    token_symbol: "LINK", direction: "sell", amount_usd: 30,
    tx_hash: "0xdef789abc012345678901234567890abcdef1234567890abcdef123456789012",
    status: "executed",
    reason: "Sold LINK because wallet 0x1234...5678 (score 87) sold it.",
    risk_assessment: "", created_at: "2025-12-08T09:15:00Z",
  },
];

export default function BotDetailPage() {
  const params = useParams();
  const id = params.id as string;
  const { isConnected } = useAuth();

  const initialBot = DUMMY_BOTS[id] || DUMMY_BOTS["1"];
  const [config, setConfig] = useState<AgentConfig>(initialBot);
  const [trades] = useState<AgentTrade[]>(DUMMY_TRADES);
  const [saving, setSaving] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const handleSaveConfig = async (data: {
    budget: number;
    max_per_trade: number;
    conditions: string[];
    consensus_threshold?: number;
    consensus_window_min?: number;
    min_score?: number;
  }) => {
    setSaving(true);
    try {
      await new Promise((r) => setTimeout(r, 800));
      setConfig((prev) => ({ ...prev, ...data }));
    } finally {
      setSaving(false);
    }
  };

  const handleStart = async () => {
    setActionLoading(true);
    try {
      await new Promise((r) => setTimeout(r, 500));
      setConfig((prev) => ({ ...prev, status: "active" }));
    } finally { setActionLoading(false); }
  };

  const handlePause = async () => {
    setActionLoading(true);
    try {
      await new Promise((r) => setTimeout(r, 500));
      setConfig((prev) => ({ ...prev, status: "paused" }));
    } finally { setActionLoading(false); }
  };

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

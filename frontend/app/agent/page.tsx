"use client";

import { useState } from "react";
import type { AgentConfig, AgentTrade } from "@/types/agent";
import { useAuth } from "@/components/providers/auth-provider";
import { AuthGuardModal } from "@/components/ui/auth-guard-modal";
import { AgentStatusCard } from "@/components/agent/agent-status-card";
import { AgentConfigForm } from "@/components/agent/agent-config-form";
import { AgentTradeHistory } from "@/components/agent/agent-trade-history";
import { Bot } from "lucide-react";

// Dummy data for visual review — will be replaced with real API calls
const DUMMY_CONFIG: AgentConfig = {
  id: 1,
  user_id: 1,
  budget: 500,
  max_per_trade: 50,
  risk_tolerance: "medium",
  min_score: 70,
  conditions: ["min_liquidity", "min_mcap"],
  status: "paused",
  agent_wallet_address: "0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18",
  total_spent: 120,
  total_trades: 4,
  created_at: "2025-12-01T10:00:00Z",
  updated_at: "2025-12-10T14:00:00Z",
};

const DUMMY_TRADES: AgentTrade[] = [
  {
    id: 1,
    agent_config_id: 1,
    source_wallet: "0x1234567890abcdef1234567890abcdef12345678",
    source_score: 87,
    token_address: "0x6982508145454ce325ddbe47a25d4ec3d2311933",
    token_symbol: "PEPE",
    direction: "buy",
    amount_usd: 45,
    tx_hash: "0xabc123def456789012345678901234567890abcdef1234567890abcdef123456",
    status: "executed",
    reason: "Bought PEPE because wallet 0x1234...5678 (score 87) bought it. Liquidity $250k, pair age 2 days.",
    risk_assessment: "Low risk — high liquidity, established token.",
    created_at: "2025-12-10T15:30:00Z",
  },
  {
    id: 2,
    agent_config_id: 1,
    source_wallet: "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
    source_score: 92,
    token_address: "0x0000000000000000000000000000000000001337",
    token_symbol: "NEWTOKEN",
    direction: "buy",
    amount_usd: 50,
    tx_hash: "",
    status: "skipped",
    reason: "Token liquidity $8k — below min_liquidity condition ($100k).",
    risk_assessment: "High risk — very low liquidity.",
    created_at: "2025-12-09T11:00:00Z",
  },
  {
    id: 3,
    agent_config_id: 1,
    source_wallet: "0x1234567890abcdef1234567890abcdef12345678",
    source_score: 87,
    token_address: "0x514910771af9ca656af840dff83e8264ecf986ca",
    token_symbol: "LINK",
    direction: "buy",
    amount_usd: 30,
    tx_hash: "0xdef789abc012345678901234567890abcdef1234567890abcdef123456789012",
    status: "executed",
    reason: "Bought LINK because wallet 0x1234...5678 (score 87) bought it. Liquidity $95M, market cap $8.5B.",
    risk_assessment: "Very low risk — top-tier token.",
    created_at: "2025-12-08T09:15:00Z",
  },
  {
    id: 4,
    agent_config_id: 1,
    source_wallet: "0x9876543210fedcba9876543210fedcba98765432",
    source_score: 65,
    token_address: "0x0000000000000000000000000000000000002222",
    token_symbol: "SCAM",
    direction: "buy",
    amount_usd: 50,
    tx_hash: "",
    status: "skipped",
    reason: "Wallet score 65 — below min_score threshold (70).",
    risk_assessment: "",
    created_at: "2025-12-07T16:45:00Z",
  },
];

export default function AgentPage() {
  const { user, isConnected } = useAuth();
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [config, setConfig] = useState<AgentConfig>(DUMMY_CONFIG);
  const [trades] = useState<AgentTrade[]>(DUMMY_TRADES);
  const [saving, setSaving] = useState(false);
  const [actionLoading, setActionLoading] = useState(false);

  const handleSaveConfig = async (data: {
    budget: number;
    max_per_trade: number;
    risk_tolerance: string;
    min_score: number;
    conditions: string[];
  }) => {
    setSaving(true);
    try {
      // TODO: Replace with real API call
      await new Promise((r) => setTimeout(r, 800));
      setConfig((prev) => ({
        ...prev,
        ...data,
        risk_tolerance: data.risk_tolerance as AgentConfig["risk_tolerance"],
      }));
    } finally {
      setSaving(false);
    }
  };

  const handleStart = async () => {
    setActionLoading(true);
    try {
      // TODO: Replace with real API call
      await new Promise((r) => setTimeout(r, 500));
      setConfig((prev) => ({ ...prev, status: "active" }));
    } finally {
      setActionLoading(false);
    }
  };

  const handlePause = async () => {
    setActionLoading(true);
    try {
      // TODO: Replace with real API call
      await new Promise((r) => setTimeout(r, 500));
      setConfig((prev) => ({ ...prev, status: "paused" }));
    } finally {
      setActionLoading(false);
    }
  };

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight">AI Trading Agent</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Set your rules. Let the agent trade for you on Base.
          </p>
        </div>

        {!isConnected ? (
          <>
            <div className="py-16 text-center">
              <Bot className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
              <p className="mb-4 text-sm text-muted-foreground">
                Connect your wallet to configure and start your AI trading agent.
              </p>
              <button
                onClick={() => setShowAuthModal(true)}
                className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
              >
                Connect Wallet
              </button>
            </div>
            <AuthGuardModal open={showAuthModal} onOpenChange={setShowAuthModal} />
          </>
        ) : (
          <div className="space-y-6">
            <AgentStatusCard
              config={config}
              onStart={handleStart}
              onPause={handlePause}
              loading={actionLoading}
            />
            <AgentConfigForm
              config={config}
              onSave={handleSaveConfig}
              saving={saving}
            />
            <AgentTradeHistory trades={trades} />
          </div>
        )}
      </div>
    </div>
  );
}

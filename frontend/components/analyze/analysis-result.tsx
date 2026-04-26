"use client";

import { useState } from "react";
import type { WalletAnalysis } from "@/types/wallet";
import { useAuth } from "@/components/providers/auth-provider";
import { AuthGuardModal } from "@/components/ui/auth-guard-modal";
import { Card, CardContent } from "@/components/ui/card";
import { ScoreRing } from "@/components/analyze/score-ring";
import { RecommendationBadge } from "@/components/analyze/recommendation-badge";
import { MetricBar } from "@/components/analyze/metric-bar";
import { TradedTokensTable } from "@/components/analyze/traded-tokens-table";
import { AiInsightCard } from "@/components/analyze/ai-insight-card";
import { ConditionsCard } from "@/components/analyze/conditions-card";
import { AlertTriangle, Eye, Check, Info } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { AttestationBadge } from "@/components/analyze/attestation-badge";
import { followWallet } from "@/api/watchlist/connector";
import { getReputation } from "@/api/reputation/connector";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription,
} from "@/components/ui/dialog";

type Props = {
  data: WalletAnalysis;
};

const METRICS = [
  {
    key: "win_rate",
    label: "Win Rate",
    tooltip:
      "Percentage of trades where the wallet made a profit. Calculated from actual PnL (buy vs sell price) using FIFO matching.",
  },
  {
    key: "profit_consistency",
    label: "Profit Consistency",
    tooltip:
      "How stable the profits are across trades. Uses standard deviation of PnL — lower spread means more consistent returns. Formula: 100 - stdDev(PnL).",
  },
  {
    key: "entry_timing",
    label: "Entry Timing",
    tooltip:
      "How early the wallet enters new tokens after launch. Based on average pair age when traded. Younger pairs = higher score (sniper behavior).",
  },
  {
    key: "token_quality",
    label: "Token Quality",
    tooltip:
      "Average market cap of tokens traded, on a logarithmic scale. Higher market cap = more established tokens. log10($10M) = score 100.",
  },
  {
    key: "trade_discipline",
    label: "Trade Discipline",
    tooltip:
      "How focused the wallet is. Ratio of unique tokens vs total transactions. Few tokens traded many times = disciplined. Many tokens traded once = scattered.",
  },
] as const;

export function AnalysisResult({ data }: Props) {
  const { user } = useAuth();
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [followed, setFollowed] = useState(false);
  const [showFollowConfirm, setShowFollowConfirm] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [attestationUID, setAttestationUID] = useState<string | undefined>();
  const [explorerURL, setExplorerURL] = useState<string | undefined>();

  // Fetch reputation (attestation) data on mount
  useState(() => {
    getReputation(data.address)
      .then((rep) => {
        setAttestationUID(rep.attestation_uid);
        setExplorerURL(rep.explorer_url);
      })
      .catch(() => {
        // No attestation yet — that's fine
      });
  });

  const handleFollow = async () => {
    if (!user) {
      setShowAuthModal(true);
      return;
    }
    setShowFollowConfirm(true);
  };

  const confirmFollow = async () => {
    setFollowLoading(true);
    try {
      await followWallet(user!.walletAddress, {
        wallet_address: data.address,
        chain: data.chain,
        recommendation: data.recommendation,
        conditions: [],
        email_notify: false,
      });
      setFollowed(true);
    } catch {
      // Silently fail
    } finally {
      setFollowLoading(false);
      setShowFollowConfirm(false);
    }
  };

  return (
    <div className="mx-auto max-w-4xl space-y-6">
      {/* Header: Score + Recommendation */}
      <Card>
        <CardContent className="flex flex-col items-center gap-6 p-6 sm:flex-row">
          <ScoreRing score={data.final_score} />
          <div className="flex-1 text-center sm:text-left">
            <div className="mb-2 flex flex-col items-center gap-2 sm:flex-row sm:justify-between">
              <div className="flex flex-col items-center gap-2 sm:flex-row">
                <RecommendationBadge recommendation={data.recommendation} />
              </div>
              {data.recommendation === "full_follow" && (
                followed ? (
                  <div className="flex items-center gap-1.5 rounded-md bg-green-500/10 px-3 py-1.5 text-sm text-green-400">
                    <Check className="h-3.5 w-3.5" />
                    Following
                  </div>
                ) : (
                  <div className="flex items-center gap-2">
                    <span className="group relative">
                      <Info className="h-4 w-4 cursor-help text-muted-foreground/50 transition-colors hover:text-muted-foreground" />
                      <span className="absolute top-full right-0 z-50 mt-2 hidden w-64 rounded-lg border bg-popover px-3 py-2 text-xs leading-relaxed text-popover-foreground shadow-lg group-hover:block">
                        Follow this wallet to get real-time notifications when it makes a trade. You&apos;ll be alerted in-app and via email so you can act fast.
                      </span>
                    </span>
                    <Button variant="outline" size="sm" className="gap-1.5" onClick={handleFollow}>
                      <Eye className="h-3.5 w-3.5" />
                      Follow Wallet
                    </Button>
                  </div>
                )
              )}
            </div>
            <p className="font-mono text-xs text-muted-foreground break-all">
              {data.address}
            </p>
            <div className="mt-2">
              <AttestationBadge attestationUID={attestationUID} explorerURL={explorerURL} />
            </div>
          </div>
        </CardContent>
      </Card>

      {/* AI Insight */}
      {data.ai_insight && (
        <AiInsightCard
          insight={data.ai_insight}
          address={data.address}
          chain={data.chain}
          tone={data.ai_insight_tone}
          prompt={data.ai_insight_prompt}
        />
      )}

      {/* Metrics */}
      <Card className="overflow-visible">
        <CardContent className="space-y-4 p-5">
          <h3 className="text-sm font-medium">Scoring Breakdown</h3>
          {METRICS.map((m, i) => (
            <MetricBar
              key={m.key}
              label={m.label}
              value={data[m.key]}
              delay={i * 100}
              tooltip={m.tooltip}
            />
          ))}
        </CardContent>
      </Card>

      {/* Conditions — interactive follow with conditions */}
      {data.conditions && data.conditions.length > 0 && (
        <ConditionsCard
          conditions={data.conditions}
          address={data.address}
          chain={data.chain}
        />
      )}

      {/* Traded Tokens */}
      {data.traded_tokens && data.traded_tokens.length > 0 && (
        <Card>
          <CardContent className="p-5">
            <h3 className="mb-4 text-sm font-medium">Recent Trades</h3>
            <TradedTokensTable tokens={data.traded_tokens} />
          </CardContent>
        </Card>
      )}

      {/* Follow Confirmation */}
      <Dialog open={showFollowConfirm} onOpenChange={setShowFollowConfirm}>
        <DialogContent className="max-w-sm">
          <DialogHeader className="items-center text-center">
            <Eye className="mb-2 h-10 w-10 text-purple-400" />
            <DialogTitle>Follow this wallet?</DialogTitle>
            <DialogDescription>
              You&apos;ll get real-time notifications when this wallet makes a trade on Base.
            </DialogDescription>
          </DialogHeader>
          <div className="flex gap-2">
            <Button variant="outline" className="flex-1" onClick={() => setShowFollowConfirm(false)}>Cancel</Button>
            <Button className="flex-1" onClick={confirmFollow} disabled={followLoading}>
              {followLoading ? "Following..." : "Follow"}
            </Button>
          </div>
        </DialogContent>
      </Dialog>

      <AuthGuardModal open={showAuthModal} onOpenChange={setShowAuthModal} />
    </div>
  );
}

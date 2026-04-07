"use client";

import type { WalletAnalysis } from "@/types/wallet";
import { Card, CardContent } from "@/components/ui/card";
import { ScoreRing } from "@/components/analyze/score-ring";
import { RecommendationBadge } from "@/components/analyze/recommendation-badge";
import { MetricBar } from "@/components/analyze/metric-bar";
import { TradedTokensTable } from "@/components/analyze/traded-tokens-table";
import { Brain, AlertTriangle } from "lucide-react";

type Props = {
  data: WalletAnalysis;
};

const METRICS = [
  { key: "win_rate", label: "Win Rate" },
  { key: "profit_consistency", label: "Profit Consistency" },
  { key: "entry_timing", label: "Entry Timing" },
  { key: "token_quality", label: "Token Quality" },
  { key: "trade_discipline", label: "Trade Discipline" },
] as const;

export function AnalysisResult({ data }: Props) {
  return (
    <div className="mx-auto max-w-4xl space-y-6">
      {/* Header: Score + Recommendation */}
      <Card>
        <CardContent className="flex flex-col items-center gap-6 p-6 sm:flex-row">
          <ScoreRing score={data.final_score} />
          <div className="flex-1 text-center sm:text-left">
            <div className="mb-2 flex flex-col items-center gap-2 sm:flex-row">
              <RecommendationBadge recommendation={data.recommendation} />
              <span className="text-sm text-muted-foreground">
                {data.total_transactions} transactions on {data.chain}
              </span>
            </div>
            <p className="font-mono text-xs text-muted-foreground break-all">
              {data.address}
            </p>
          </div>
        </CardContent>
      </Card>

      {/* AI Insight */}
      {data.ai_insight && (
        <Card>
          <CardContent className="flex gap-3 p-5">
            <Brain className="mt-0.5 h-5 w-5 shrink-0 text-purple-400" />
            <p className="text-sm leading-relaxed text-muted-foreground">
              {data.ai_insight}
            </p>
          </CardContent>
        </Card>
      )}

      {/* Metrics */}
      <Card>
        <CardContent className="space-y-4 p-5">
          <h3 className="text-sm font-medium">Scoring Breakdown</h3>
          {METRICS.map((m, i) => (
            <MetricBar
              key={m.key}
              label={m.label}
              value={data[m.key]}
              delay={i * 100}
            />
          ))}
          <div className="flex items-center gap-2 rounded-lg bg-muted/50 px-3 py-2 text-xs text-muted-foreground">
            <AlertTriangle className="h-3.5 w-3.5" />
            Risk Exposure: {data.risk_exposure.toFixed(1)}% (informational)
          </div>
        </CardContent>
      </Card>

      {/* Conditions (conditional_follow only) */}
      {data.conditions && data.conditions.length > 0 && (
        <Card>
          <CardContent className="space-y-3 p-5">
            <h3 className="text-sm font-medium">Suggested Conditions</h3>
            <div className="flex flex-wrap gap-2">
              {data.conditions.map((c) => (
                <span
                  key={c.id}
                  className="rounded-full border border-yellow-500/30 bg-yellow-500/10 px-3 py-1 text-xs text-yellow-400"
                >
                  {c.label}
                </span>
              ))}
            </div>
          </CardContent>
        </Card>
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
    </div>
  );
}

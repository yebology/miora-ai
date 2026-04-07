"use client";

import { useParams } from "next/navigation";
import { useState, useEffect } from "react";
import type { WalletAnalysis } from "@/types/wallet";
import type { Notification } from "@/types/watchlist";
import { DUMMY_ANALYSIS } from "@/constants/dummy";
import { DUMMY_NOTIFICATIONS } from "@/constants/dummy-watchlist";
import { AnalysisResult } from "@/components/analyze/analysis-result";
import { NotificationItem } from "@/components/watchlist/notification-item";
import { Button } from "@/components/ui/button";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import Link from "next/link";
import {
  ArrowLeft,
  RefreshCw,
  Bell,
  Loader2,
  CheckCircle,
  AlertTriangle,
} from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";

export default function WatchlistDetailPage() {
  const params = useParams();
  const chain = params.chain as string;
  const address = params.address as string;

  const [analysis, setAnalysis] = useState<WalletAnalysis | null>(null);
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [reanalyzing, setReanalyzing] = useState(false);
  const [showSuccess, setShowSuccess] = useState(false);
  const [showConfirm, setShowConfirm] = useState(false);
  const [tab, setTab] = useState<"analysis" | "activity">("analysis");

  useEffect(() => {
    const load = async () => {
      await new Promise((r) => setTimeout(r, 500));
      setAnalysis({ ...DUMMY_ANALYSIS, address, chain });
      setNotifications(
        DUMMY_NOTIFICATIONS.filter(
          (n) => n.wallet_address === address && n.chain === chain
        )
      );
      setLoading(false);
    };
    load();
  }, [address, chain]);

  const handleReanalyze = async () => {
    setShowConfirm(false);
    setReanalyzing(true);
    try {
      await new Promise((r) => setTimeout(r, 2000));
      setAnalysis({
        ...DUMMY_ANALYSIS,
        address,
        chain,
        final_score: Math.round(Math.random() * 30 + 50),
      });
      setShowSuccess(true);
    } finally {
      setReanalyzing(false);
    }
  };

  if (loading) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <div className="flex flex-col items-center gap-3">
          <div className="h-8 w-8 animate-spin rounded-full border-2 border-muted border-t-primary" />
          <p className="text-sm text-muted-foreground">Loading wallet data...</p>
        </div>
      </div>
    );
  }

  if (!analysis) {
    return (
      <div className="flex flex-1 items-center justify-center py-24">
        <div className="text-center">
          <p className="text-sm text-muted-foreground">Wallet not found.</p>
          <Link
            href="/watchlist"
            className={cn(buttonVariants({ variant: "outline", size: "sm" }), "mt-4")}
          >
            Back to Watchlist
          </Link>
        </div>
      </div>
    );
  }

  const hasConditions =
    analysis.recommendation === "conditional_follow" &&
    analysis.conditions &&
    analysis.conditions.length > 0;

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-4xl">
        <div className="mb-6 flex items-center justify-between">
          <Link
            href="/watchlist"
            className="flex items-center gap-1.5 text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <ArrowLeft className="h-4 w-4" />
            Back to Watchlist
          </Link>
          <Button
            variant="outline"
            size="sm"
            className="gap-1.5"
            onClick={() => setShowConfirm(true)}
            disabled={reanalyzing}
          >
            {reanalyzing ? (
              <Loader2 className="h-3.5 w-3.5 animate-spin" />
            ) : (
              <RefreshCw className="h-3.5 w-3.5" />
            )}
            {reanalyzing ? "Analyzing..." : "Re-analyze"}
          </Button>
        </div>

        <div className="mb-6 flex gap-1 rounded-lg bg-muted/50 p-1">
          <button
            onClick={() => setTab("analysis")}
            className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              tab === "analysis"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            Stored Analysis
          </button>
          <button
            onClick={() => setTab("activity")}
            className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              tab === "activity"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Bell className="h-4 w-4" />
            Activity ({notifications.length})
          </button>
        </div>

        {tab === "analysis" && <AnalysisResult data={analysis} />}

        {tab === "activity" && (
          <div className="space-y-2">
            {notifications.length === 0 ? (
              <div className="py-16 text-center">
                <Bell className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
                <p className="text-sm text-muted-foreground">
                  No trade activity detected yet for this wallet.
                </p>
              </div>
            ) : (
              notifications.map((n) => (
                <NotificationItem key={n.id} notification={n} />
              ))
            )}
          </div>
        )}

        {/* Confirm Re-analyze */}
        <Dialog open={showConfirm} onOpenChange={setShowConfirm}>
          <DialogContent className="max-w-sm">
            <DialogHeader className="items-center text-center">
              <AlertTriangle className="mb-2 h-10 w-10 text-yellow-400" />
              <DialogTitle>Re-analyze this wallet?</DialogTitle>
              <DialogDescription>
                This will fetch the latest on-chain data and overwrite the
                current scoring and recommendations.
                {hasConditions && (
                  <span className="mt-2 block text-yellow-400">
                    Your selected follow conditions will need to be re-selected
                    if the recommendation changes.
                  </span>
                )}
              </DialogDescription>
            </DialogHeader>
            <div className="flex gap-2">
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => setShowConfirm(false)}
              >
                Cancel
              </Button>
              <Button className="flex-1" onClick={handleReanalyze}>
                Re-analyze
              </Button>
            </div>
          </DialogContent>
        </Dialog>

        {/* Success Modal */}
        <Dialog open={showSuccess} onOpenChange={setShowSuccess}>
          <DialogContent className="max-w-sm text-center">
            <DialogHeader className="items-center">
              <CheckCircle className="mb-2 h-10 w-10 text-green-400" />
              <DialogTitle>Analysis Updated</DialogTitle>
              <DialogDescription>
                Wallet has been re-analyzed with the latest on-chain data. Score
                and recommendations have been refreshed.
              </DialogDescription>
            </DialogHeader>
            <Button className="w-full" onClick={() => setShowSuccess(false)}>
              Got it
            </Button>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}

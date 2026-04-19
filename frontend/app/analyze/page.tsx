"use client";

import { useState } from "react";
import type { WalletAnalysis } from "@/api/wallet/validation";
import { analyzeWallet, getWallet } from "@/api/wallet/connector";
import { AnalyzeForm } from "@/components/analyze/analyze-form";
import { AnalysisResult } from "@/components/analyze/analysis-result";
import { Button } from "@/components/ui/button";
import { AlertTriangle } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";

export default function AnalyzePage() {
  const [result, setResult] = useState<WalletAnalysis | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showExistsConfirm, setShowExistsConfirm] = useState(false);
  const [pendingRequest, setPendingRequest] = useState<{
    address: string;
    chain: string;
  } | null>(null);

  const checkAndAnalyze = async (address: string, chain: string) => {
    setError(null);
    setResult(null);

    try {
      const existing = await getWallet(address);
      if (existing) {
        setPendingRequest({ address, chain });
        setShowExistsConfirm(true);
        return;
      }
    } catch {
      // Wallet not found — proceed to analyze
    }

    await doAnalyze(address, chain);
  };

  const doAnalyze = async (address: string, chain: string) => {
    setLoading(true);
    setResult(null);
    setError(null);
    setShowExistsConfirm(false);

    try {
      const data = await analyzeWallet(address, chain);
      setResult(data);
    } catch (err) {
      setError(
        err instanceof Error
          ? err.message
          : "Failed to analyze wallet. Please try again."
      );
    } finally {
      setLoading(false);
    }
  };

  const handleConfirmReanalyze = () => {
    if (pendingRequest) {
      doAnalyze(pendingRequest.address, pendingRequest.chain);
    }
  };

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-4xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight">
            Wallet Analysis
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Paste any wallet address to get AI-powered scoring and
            recommendations.
          </p>
        </div>

        <AnalyzeForm
          onAnalyze={checkAndAnalyze}
          loading={loading}
          error={error}
        />

        {loading && (
          <div className="mt-12 flex flex-col items-center gap-3">
            <div className="h-8 w-8 animate-spin rounded-full border-2 border-muted border-t-primary" />
            <p className="text-sm text-muted-foreground">
              Analyzing wallet...
            </p>
          </div>
        )}

        {result && (
          <div className="mt-8">
            <AnalysisResult data={result} />
          </div>
        )}

        {/* Wallet Already Exists Confirmation */}
        <Dialog open={showExistsConfirm} onOpenChange={setShowExistsConfirm}>
          <DialogContent className="max-w-sm">
            <DialogHeader className="items-center text-center">
              <AlertTriangle className="mb-2 h-10 w-10 text-yellow-400" />
              <DialogTitle>Wallet already analyzed</DialogTitle>
              <DialogDescription>
                This wallet has been analyzed before. Re-analyzing will fetch
                the latest data and overwrite the previous scoring,
                recommendations, and AI insights.
              </DialogDescription>
            </DialogHeader>
            <div className="flex gap-2">
              <Button
                variant="outline"
                className="flex-1"
                onClick={() => setShowExistsConfirm(false)}
              >
                Cancel
              </Button>
              <Button className="flex-1" onClick={handleConfirmReanalyze}>
                Re-analyze
              </Button>
            </div>
          </DialogContent>
        </Dialog>
      </div>
    </div>
  );
}

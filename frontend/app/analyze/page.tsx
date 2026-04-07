"use client";

import { useState } from "react";
import type { WalletAnalysis } from "@/types/wallet";
import { AnalyzeForm } from "@/components/analyze/analyze-form";
import { AnalysisResult } from "@/components/analyze/analysis-result";
import { DUMMY_ANALYSIS } from "@/constants/dummy";

export default function AnalyzePage() {
  const [result, setResult] = useState<WalletAnalysis | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleAnalyze = async (address: string, chain: string) => {
    setLoading(true);
    setResult(null);
    setError(null);

    try {
      // TODO: Replace with real API call
      // const res = await fetch(`${API_URL}/api/wallets/analyze`, {
      //   method: "POST",
      //   headers: { "Content-Type": "application/json" },
      //   body: JSON.stringify({ address, chain }),
      // });
      // const json = await res.json();
      // if (json.status === "error") throw new Error(json.message);
      // setResult(json.data);

      // Simulate loading with dummy data
      await new Promise((r) => setTimeout(r, 1500));

      // Type "error" to simulate error state (for testing)
      if (address.toLowerCase() === "error") {
        throw new Error("Wallet not found. Please check the address and try again.");
      }

      setResult({ ...DUMMY_ANALYSIS, address, chain });
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
          onAnalyze={handleAnalyze}
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
      </div>
    </div>
  );
}

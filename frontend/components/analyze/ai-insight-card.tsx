"use client";

import { useState } from "react";
import { Brain, RefreshCw, ChevronDown, Send } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { regenerateInsight } from "@/api/wallet/connector";

const TONES = [
  { value: "simple", label: "Simple", emoji: "💬" },
  { value: "eli5", label: "ELI5", emoji: "🧒" },
  { value: "custom", label: "Custom", emoji: "✏️" },
] as const;

type Props = {
  insight: string;
  address: string;
  chain: string;
  tone?: string;
  prompt?: string;
};

export function AiInsightCard({ insight, address, chain, tone, prompt }: Props) {
  const [currentInsight, setCurrentInsight] = useState(insight);
  const [activeTone, setActiveTone] = useState(tone || "simple");
  const [loading, setLoading] = useState(false);
  const [showTones, setShowTones] = useState(false);
  const [showCustomInput, setShowCustomInput] = useState(false);
  const [customPrompt, setCustomPrompt] = useState("");

  const handleRegenerate = async (tone: string, prompt?: string) => {
    setLoading(true);
    setShowTones(false);
    setActiveTone(tone);

    try {
      const res = await regenerateInsight(
        address,
        chain,
        tone as "simple" | "eli5" | "custom",
        prompt,
      );
      setCurrentInsight(res.ai_insight);
      if (tone === "custom") setShowCustomInput(false);
    } catch {
      // Keep current insight on error
    } finally {
      setLoading(false);
    }
  };

  const handleCustomSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!customPrompt.trim()) return;
    handleRegenerate("custom", customPrompt.trim());
  };

  return (
    <Card>
      <CardContent className="p-5">
        <div className="mb-3 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <Brain className="h-5 w-5 text-purple-400" />
            <span className="text-sm font-medium">AI Insight</span>
          </div>

          <div className="relative">
            <Button
              variant="ghost"
              size="sm"
              className="gap-1.5 text-xs text-muted-foreground"
              onClick={() => {
                setShowTones(!showTones);
                setShowCustomInput(false);
              }}
              disabled={loading}
            >
              {loading ? (
                <RefreshCw className="h-3.5 w-3.5 animate-spin" />
              ) : (
                <RefreshCw className="h-3.5 w-3.5" />
              )}
              {TONES.find((t) => t.value === activeTone)?.label}
              <ChevronDown className="h-3 w-3" />
            </Button>

            {showTones && (
              <div className="absolute right-0 top-full z-50 mt-1 w-40 rounded-lg border bg-popover p-1 shadow-lg">
                {TONES.map((tone) => (
                  <button
                    key={tone.value}
                    onClick={() => {
                      if (tone.value === "custom") {
                        setShowTones(false);
                        setShowCustomInput(true);
                      } else {
                        handleRegenerate(tone.value);
                      }
                    }}
                    className={`flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition-colors hover:bg-muted ${
                      activeTone === tone.value
                        ? "bg-muted text-foreground"
                        : "text-muted-foreground"
                    }`}
                  >
                    <span>{tone.emoji}</span>
                    {tone.label}
                  </button>
                ))}
              </div>
            )}
          </div>
        </div>

        <p
          className={`text-sm leading-relaxed text-muted-foreground transition-opacity duration-300 ${
            loading ? "opacity-50" : "opacity-100"
          }`}
        >
          {currentInsight}
        </p>

        {/* Custom prompt input */}
        {showCustomInput && (
          <form
            onSubmit={handleCustomSubmit}
            className="mt-3 flex gap-2"
          >
            <Input
              placeholder="Ask anything about this wallet..."
              value={customPrompt}
              onChange={(e) => setCustomPrompt(e.target.value)}
              maxLength={200}
              className="h-9 text-sm"
              autoFocus
            />
            <Button
              type="submit"
              size="sm"
              disabled={!customPrompt.trim() || loading}
              className="h-9 gap-1.5"
            >
              <Send className="h-3.5 w-3.5" />
              Ask
            </Button>
          </form>
        )}
      </CardContent>
    </Card>
  );
}

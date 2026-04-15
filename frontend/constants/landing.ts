import {
  Brain,
  Shield,
  Bell,
  Bot,
  Zap,
  Search,
  BarChart3,
  MousePointerClick,
  Fingerprint,
  DollarSign,
  type LucideIcon,
} from "lucide-react";

export type Feature = {
  icon: LucideIcon;
  title: string;
  description: string;
};

export type Chain = {
  name: string;
  logo: string;
  color: string;
};

export type Step = {
  step: string;
  title: string;
  description: string;
  icon: LucideIcon;
};

export const FEATURES: Feature[] = [
  {
    icon: Fingerprint,
    title: "On-chain Reputation (EAS)",
    description:
      "Trading scores published on-chain via Ethereum Attestation Service on Base. Verifiable by any protocol, agent, or dApp.",
  },
  {
    icon: Brain,
    title: "AI-Powered Insights",
    description:
      "Understand any wallet's trading behavior in plain language. No charts, no jargon — just clear, actionable analysis.",
  },
  {
    icon: Shield,
    title: "Smart Recommendations",
    description:
      "Get Full Follow, Conditional Follow, or Avoid ratings with AI-generated conditions for every wallet.",
  },
  {
    icon: Bell,
    title: "Real-Time Alerts",
    description:
      "Follow wallets and get notified instantly when they trade, filtered by your custom conditions.",
  },
  {
    icon: Bot,
    title: "AI Trading Agent",
    description:
      "Autonomous agent monitors top wallets, evaluates trades, and executes swaps on Base via Coinbase AgentKit.",
  },
  {
    icon: DollarSign,
    title: "x402 Reputation API",
    description:
      "Other protocols and AI agents can query Miora scores via USDC micropayments. No API keys needed.",
  },
];

export const CHAINS: Chain[] = [
  { name: "Base", logo: "/chains/base.svg", color: "#0052FF" },
];

export const STEPS: Step[] = [
  {
    step: "01",
    title: "Analyze",
    description: "Paste any wallet address on Base to get a multi-factor trading score.",
    icon: Search,
  },
  {
    step: "02",
    title: "Follow or Agent",
    description:
      "Follow the wallet with smart conditions, or let the AI agent trade for you automatically.",
    icon: BarChart3,
  },
  {
    step: "03",
    title: "Verified On-chain",
    description:
      "Every score is published as an EAS attestation on Base — verifiable by anyone.",
    icon: MousePointerClick,
  },
];

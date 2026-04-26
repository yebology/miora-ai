import {
  Brain,
  Shield,
  Bell,
  Bot,
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
      "Every trading score is published as an EAS attestation on Base. Queryable by any protocol, agent, or dApp — a composable reputation primitive.",
  },
  {
    icon: Brain,
    title: "AI-Powered Insights",
    description:
      "Gemini AI translates complex scoring data into plain-language analysis. Ask in simple, ELI5, or custom tone — no charts or jargon needed.",
  },
  {
    icon: Shield,
    title: "Dynamic Conditions",
    description:
      "Full Follow, Conditional Follow, or Avoid — with filter thresholds computed from the wallet's own trading data, not hardcoded numbers.",
  },
  {
    icon: Bell,
    title: "Smart Alerts + AI Risk",
    description:
      "Get notified when followed wallets trade. Each alert includes an AI risk assessment evaluating liquidity, market cap, and pair age before you act.",
  },
  {
    icon: Bot,
    title: "AI Trading Bots (AgentKit)",
    description:
      "Wallet Bot copies one trader's buys and sells. Consensus Bot trades when 3+ top wallets agree. Each bot has its own Agentic Wallet on Base.",
  },
  {
    icon: DollarSign,
    title: "Auto Profit Transfer",
    description:
      "When a bot sells, proceeds are automatically transferred on-chain to your connected wallet. No manual withdrawal needed.",
  },
];

export const CHAINS: Chain[] = [
  { name: "Base", logo: "/chains/base.svg", color: "#0052FF" },
];

export const STEPS: Step[] = [
  {
    step: "01",
    title: "Analyze",
    description: "Paste any wallet address on Base. Get a multi-factor score, recommendation, and AI insight in seconds.",
    icon: Search,
  },
  {
    step: "02",
    title: "Follow or Automate",
    description:
      "Follow with smart conditions and get AI-assessed alerts, or create a bot that trades for you automatically.",
    icon: BarChart3,
  },
  {
    step: "03",
    title: "Verified On-chain",
    description:
      "Every score is published as an EAS attestation on Base — verifiable and queryable by anyone.",
    icon: MousePointerClick,
  },
];

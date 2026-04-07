import {
  Brain,
  ArrowLeftRight,
  Shield,
  Bell,
  Globe,
  Zap,
  Search,
  BarChart3,
  MousePointerClick,
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
  color: string; // tailwind shadow color for hover glow
};

export type Step = {
  step: string;
  title: string;
  description: string;
  icon: LucideIcon;
};

export const FEATURES: Feature[] = [
  {
    icon: Brain,
    title: "AI-Powered Insights",
    description:
      "Understand any wallet's trading behavior in plain language. No charts, no jargon — just clear, actionable analysis.",
  },
  {
    icon: ArrowLeftRight,
    title: "DEX Aggregator",
    description:
      "Swap tokens across Jupiter (Solana) and 1inch (EVM) with best route discovery, all in one place.",
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
    icon: Globe,
    title: "Multi-Chain Support",
    description:
      "Ethereum, Arbitrum, Optimism, Base, Polygon, and Solana — all supported out of the box.",
  },
  {
    icon: Zap,
    title: "Beginner Friendly",
    description:
      "Built for everyone, not just advanced traders. AI translates complex data into simple decisions.",
  },
];

export const CHAINS: Chain[] = [
  { name: "Ethereum", logo: "/chains/ethereum.svg", color: "#627EEA" },
  { name: "Arbitrum", logo: "/chains/arbitrum.svg", color: "#28A0F0" },
  { name: "Optimism", logo: "/chains/optimism.svg", color: "#FF0420" },
  { name: "Base", logo: "/chains/base.svg", color: "#0052FF" },
  { name: "Polygon", logo: "/chains/polygon.svg", color: "#8247E5" },
  { name: "Solana", logo: "/chains/solana.svg", color: "#14F195" },
];

export const STEPS: Step[] = [
  {
    step: "01",
    title: "Analyze",
    description: "Paste any wallet address and select a chain to start.",
    icon: Search,
  },
  {
    step: "02",
    title: "Understand",
    description:
      "Get a multi-factor score with AI-powered insights in plain language.",
    icon: BarChart3,
  },
  {
    step: "03",
    title: "Decide",
    description:
      "Follow the wallet, set alert conditions, or swap tokens directly.",
    icon: MousePointerClick,
  },
];

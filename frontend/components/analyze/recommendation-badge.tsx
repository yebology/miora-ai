import { cn } from "@/lib/utils";

const CONFIG = {
  full_follow: {
    label: "Full Follow",
    emoji: "🟢",
    className: "border-green-500/30 bg-green-500/10 text-green-400",
  },
  conditional_follow: {
    label: "Conditional Follow",
    emoji: "🟡",
    className: "border-yellow-500/30 bg-yellow-500/10 text-yellow-400",
  },
  avoid: {
    label: "Avoid",
    emoji: "🔴",
    className: "border-red-500/30 bg-red-500/10 text-red-400",
  },
} as const;

type Props = {
  recommendation: keyof typeof CONFIG;
};

export function RecommendationBadge({ recommendation }: Props) {
  const config = CONFIG[recommendation];

  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full border px-3 py-1 text-sm font-medium",
        config.className
      )}
    >
      {config.emoji} {config.label}
    </span>
  );
}

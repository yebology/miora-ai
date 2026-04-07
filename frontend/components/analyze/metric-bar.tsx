"use client";

import { useEffect, useState } from "react";

type Props = {
  label: string;
  value: number;
  max?: number;
  delay?: number;
};

export function MetricBar({ label, value, max = 100, delay = 0 }: Props) {
  const [width, setWidth] = useState(0);
  const percent = Math.min((value / max) * 100, 100);

  useEffect(() => {
    const timer = setTimeout(() => setWidth(percent), 200 + delay);
    return () => clearTimeout(timer);
  }, [percent, delay]);

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-sm">
        <span className="text-muted-foreground">{label}</span>
        <span className="font-mono font-medium">{value.toFixed(1)}</span>
      </div>
      <div className="h-2 overflow-hidden rounded-full bg-muted">
        <div
          className="h-full rounded-full bg-primary transition-all duration-700 ease-out"
          style={{ width: `${width}%` }}
        />
      </div>
    </div>
  );
}

"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { Info } from "lucide-react";

type Props = {
  label: string;
  value: number;
  max?: number;
  delay?: number;
  tooltip?: string;
};

export function MetricBar({ label, value, max = 100, delay = 0, tooltip }: Props) {
  const [width, setWidth] = useState(0);
  const [showTooltip, setShowTooltip] = useState(false);
  const [position, setPosition] = useState<"top" | "bottom">("top");
  const iconRef = useRef<HTMLSpanElement>(null);
  const percent = Math.min((value / max) * 100, 100);

  useEffect(() => {
    const timer = setTimeout(() => setWidth(percent), 200 + delay);
    return () => clearTimeout(timer);
  }, [percent, delay]);

  const handleMouseEnter = useCallback(() => {
    if (iconRef.current) {
      const rect = iconRef.current.getBoundingClientRect();
      // If less than 120px from top of viewport, show below; otherwise show above
      setPosition(rect.top < 120 ? "bottom" : "top");
    }
    setShowTooltip(true);
  }, []);

  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-sm">
        <span className="flex items-center gap-1.5 text-muted-foreground">
          {label}
          {tooltip && (
            <span
              ref={iconRef}
              className="relative"
              onMouseEnter={handleMouseEnter}
              onMouseLeave={() => setShowTooltip(false)}
            >
              <Info className="h-3.5 w-3.5 cursor-help text-muted-foreground/50 transition-colors hover:text-muted-foreground" />
              {showTooltip && (
                <span
                  className={`absolute left-0 z-50 w-64 rounded-lg border bg-popover px-3 py-2 text-xs leading-relaxed text-popover-foreground shadow-lg ${
                    position === "top" ? "bottom-full mb-2" : "top-full mt-2"
                  }`}
                >
                  {tooltip}
                </span>
              )}
            </span>
          )}
        </span>
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

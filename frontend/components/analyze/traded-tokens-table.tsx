"use client";

import { useState } from "react";
import type { TradedToken } from "@/types/wallet";
import { cn } from "@/lib/utils";
import { ArrowUpDown } from "lucide-react";

type Props = {
  tokens: TradedToken[];
};

type SortKey = "symbol" | "pnl_percent" | "buy_price" | "exit_price";
type SortDir = "asc" | "desc";

function TradeSummary({ tokens }: Props) {
  const total = tokens.length;
  const profitable = tokens.filter((t) => t.pnl_percent > 0).length;
  const loss = tokens.filter((t) => t.pnl_percent < 0).length;
  const avgPnl = tokens.reduce((sum, t) => sum + t.pnl_percent, 0) / total;
  const realized = tokens.filter((t) => t.status === "realized").length;

  return (
    <div className="mb-4 grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div className="rounded-lg bg-muted/50 px-3 py-2">
        <p className="text-xs text-muted-foreground">Tokens Traded</p>
        <p className="text-lg font-semibold">{total}</p>
      </div>
      <div className="rounded-lg bg-muted/50 px-3 py-2">
        <p className="text-xs text-muted-foreground">Avg PnL</p>
        <p
          className={cn(
            "text-lg font-semibold font-mono",
            avgPnl >= 0 ? "text-green-400" : "text-red-400"
          )}
        >
          {avgPnl >= 0 ? "+" : ""}{avgPnl.toFixed(1)}%
        </p>
      </div>
      <div className="rounded-lg bg-muted/50 px-3 py-2">
        <p className="text-xs text-muted-foreground">Win / Loss</p>
        <p className="text-lg font-semibold">
          <span className="text-green-400">{profitable}</span>
          <span className="text-muted-foreground"> / </span>
          <span className="text-red-400">{loss}</span>
        </p>
      </div>
      <div className="rounded-lg bg-muted/50 px-3 py-2">
        <p className="text-xs text-muted-foreground">Realized</p>
        <p className="text-lg font-semibold">
          {realized}
          <span className="text-sm font-normal text-muted-foreground"> / {total}</span>
        </p>
      </div>
    </div>
  );
}

const COLUMNS: { key: SortKey; label: string }[] = [
  { key: "symbol", label: "Token" },
  { key: "pnl_percent", label: "PnL (%)" },
  { key: "buy_price", label: "Buy (USD)" },
  { key: "exit_price", label: "Exit (USD)" },
];

export function TradedTokensTable({ tokens }: Props) {
  const [sortKey, setSortKey] = useState<SortKey>("pnl_percent");
  const [sortDir, setSortDir] = useState<SortDir>("desc");

  const handleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir(sortDir === "asc" ? "desc" : "asc");
    } else {
      setSortKey(key);
      setSortDir("desc");
    }
  };

  const sorted = [...tokens].sort((a, b) => {
    const aVal = a[sortKey];
    const bVal = b[sortKey];
    if (typeof aVal === "string" && typeof bVal === "string") {
      return sortDir === "asc" ? aVal.localeCompare(bVal) : bVal.localeCompare(aVal);
    }
    return sortDir === "asc"
      ? (aVal as number) - (bVal as number)
      : (bVal as number) - (aVal as number);
  });

  return (
    <div>
      <TradeSummary tokens={tokens} />

      <div className="max-h-80 overflow-y-auto rounded-lg border">
        <table className="w-full text-sm">
          <thead className="sticky top-0 bg-card">
            <tr className="border-b text-left text-muted-foreground">
              {COLUMNS.map((col) => (
                <th
                  key={col.key}
                  className="cursor-pointer px-3 py-2.5 font-medium select-none transition-colors hover:text-foreground"
                  onClick={() => handleSort(col.key)}
                >
                  <span className="flex items-center gap-1">
                    {col.label}
                    <ArrowUpDown
                      className={cn(
                        "h-3 w-3",
                        sortKey === col.key
                          ? "text-foreground"
                          : "text-muted-foreground/30"
                      )}
                    />
                  </span>
                </th>
              ))}
              <th className="px-3 py-2.5 font-medium">Status</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((token, i) => (
              <tr
                key={`${token.contract_address}-${i}`}
                className="border-b border-border/50 last:border-0"
              >
                <td className="px-3 py-2.5 font-mono font-medium">
                  {token.symbol}
                </td>
                <td
                  className={cn(
                    "px-3 py-2.5 font-mono",
                    token.pnl_percent >= 0 ? "text-green-400" : "text-red-400"
                  )}
                >
                  {token.pnl_percent >= 0 ? "+" : ""}
                  {token.pnl_percent.toFixed(1)}%
                </td>
                <td className="px-3 py-2.5 font-mono text-muted-foreground">
                  ${formatPrice(token.buy_price)}
                </td>
                <td className="px-3 py-2.5 font-mono text-muted-foreground">
                  ${formatPrice(token.exit_price)}
                </td>
                <td className="px-3 py-2.5">
                  <span
                    className={cn(
                      "rounded-full px-2 py-0.5 text-xs",
                      token.status === "realized"
                        ? "bg-muted text-muted-foreground"
                        : "bg-blue-500/10 text-blue-400"
                    )}
                  >
                    {token.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function formatPrice(price: number): string {
  if (price === 0) return "0.00";
  if (price < 0.01) {
    const str = price.toFixed(20);
    const match = str.match(/^0\.(0*)/);
    const zeros = match ? match[1].length : 0;
    return price.toFixed(zeros + 2);
  }
  if (price < 1) return price.toFixed(4);
  if (price < 1000) return price.toFixed(2);
  return price.toLocaleString(undefined, { maximumFractionDigits: 2 });
}

import type { TradedToken } from "@/types/wallet";
import { cn } from "@/lib/utils";

type Props = {
  tokens: TradedToken[];
};

export function TradedTokensTable({ tokens }: Props) {
  return (
    <div className="max-h-80 overflow-y-auto rounded-lg border">
      <table className="w-full text-sm">
        <thead className="sticky top-0 bg-card">
          <tr className="border-b text-left text-muted-foreground">
            <th className="px-3 py-2.5 font-medium">Token</th>
            <th className="px-3 py-2.5 font-medium">PnL (%)</th>
            <th className="px-3 py-2.5 font-medium">Buy (USD)</th>
            <th className="px-3 py-2.5 font-medium">Exit (USD)</th>
            <th className="px-3 py-2.5 font-medium">Status</th>
          </tr>
        </thead>
        <tbody>
          {tokens.map((token, i) => (
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
  );
}

function formatPrice(price: number): string {
  if (price === 0) return "0.00";
  if (price < 0.01) {
    // Show full decimals for tiny prices (memecoins)
    // Count leading zeros after decimal, then show 2 significant digits
    const str = price.toFixed(20);
    const match = str.match(/^0\.(0*)/);
    const zeros = match ? match[1].length : 0;
    return price.toFixed(zeros + 2);
  }
  if (price < 1) return price.toFixed(4);
  if (price < 1000) return price.toFixed(2);
  return price.toLocaleString(undefined, { maximumFractionDigits: 2 });
}

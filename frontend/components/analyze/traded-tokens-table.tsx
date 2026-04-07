import type { TradedToken } from "@/types/wallet";
import { cn } from "@/lib/utils";

type Props = {
  tokens: TradedToken[];
};

export function TradedTokensTable({ tokens }: Props) {
  return (
    <div className="overflow-x-auto">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b text-left text-muted-foreground">
            <th className="pb-2 font-medium">Token</th>
            <th className="pb-2 font-medium">PnL</th>
            <th className="pb-2 font-medium">Buy</th>
            <th className="pb-2 font-medium">Exit</th>
            <th className="pb-2 font-medium">Status</th>
          </tr>
        </thead>
        <tbody>
          {tokens.map((token) => (
            <tr key={token.contract_address} className="border-b border-border/50">
              <td className="py-2.5 font-mono font-medium">{token.symbol}</td>
              <td
                className={cn(
                  "py-2.5 font-mono",
                  token.pnl_percent >= 0 ? "text-green-400" : "text-red-400"
                )}
              >
                {token.pnl_percent >= 0 ? "+" : ""}
                {token.pnl_percent.toFixed(1)}%
              </td>
              <td className="py-2.5 font-mono text-muted-foreground">
                {formatPrice(token.buy_price)}
              </td>
              <td className="py-2.5 font-mono text-muted-foreground">
                {formatPrice(token.exit_price)}
              </td>
              <td className="py-2.5">
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
  if (price < 0.001) return price.toExponential(2);
  if (price < 1) return price.toFixed(6);
  return price.toLocaleString(undefined, { maximumFractionDigits: 2 });
}

"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Wallet, Loader2, CheckCircle, ExternalLink } from "lucide-react";
import {
  Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription,
} from "@/components/ui/dialog";
import { useAppKitAccount } from "@reown/appkit/react";
import { useSendTransaction, useWaitForTransactionReceipt } from "wagmi";
import { encodeFunctionData, parseUnits } from "viem";

const MUSDT_ADDRESS = process.env.NEXT_PUBLIC_MUSDT_ADDRESS as `0x${string}`;

// ERC-20 transfer function signature
const ERC20_TRANSFER_ABI = [
  {
    name: "transfer",
    type: "function",
    inputs: [
      { name: "to", type: "address" },
      { name: "amount", type: "uint256" },
    ],
    outputs: [{ name: "", type: "bool" }],
  },
] as const;

type Props = {
  agentWalletAddress: string;
  onDepositSuccess?: () => void;
};

export function DepositButton({ agentWalletAddress, onDepositSuccess }: Props) {
  const { isConnected } = useAppKitAccount();
  const [showDialog, setShowDialog] = useState(false);
  const [amount, setAmount] = useState("");
  const [txHash, setTxHash] = useState<`0x${string}` | undefined>();

  const { sendTransaction, isPending: isSending } = useSendTransaction();
  const { isLoading: isConfirming, isSuccess } = useWaitForTransactionReceipt({
    hash: txHash,
  });

  const handleDeposit = () => {
    if (!amount || !agentWalletAddress || !MUSDT_ADDRESS) return;

    const data = encodeFunctionData({
      abi: ERC20_TRANSFER_ABI,
      functionName: "transfer",
      args: [agentWalletAddress as `0x${string}`, parseUnits(amount, 6)],
    });

    sendTransaction(
      {
        to: MUSDT_ADDRESS,
        data,
      },
      {
        onSuccess: (hash) => {
          setTxHash(hash);
        },
      },
    );
  };

  // Auto-callback on success
  if (isSuccess && onDepositSuccess) {
    onDepositSuccess();
  }

  if (!agentWalletAddress) {
    return (
      <div className="rounded-lg bg-yellow-500/10 px-3 py-2 text-xs text-yellow-400">
        Agent wallet not ready. Start the agent sidecar first.
      </div>
    );
  }

  return (
    <>
      <Button
        variant="outline"
        className="w-full gap-1.5"
        onClick={() => setShowDialog(true)}
        disabled={!isConnected}
      >
        <Wallet className="h-3.5 w-3.5" />
        Deposit USDT
      </Button>

      <Dialog open={showDialog} onOpenChange={setShowDialog}>
        <DialogContent className="max-w-sm">
          <DialogHeader className="items-center text-center">
            <Wallet className="mb-2 h-10 w-10 text-purple-400" />
            <DialogTitle>Deposit USDT</DialogTitle>
            <DialogDescription>
              Transfer USDT to the bot&apos;s Agentic Wallet so it can trade.
            </DialogDescription>
          </DialogHeader>

          <div className="rounded-lg bg-muted/30 px-3 py-2">
            <p className="text-xs text-muted-foreground">Agent Wallet</p>
            <p className="font-mono text-xs break-all">{agentWalletAddress}</p>
          </div>

          {isSuccess ? (
            <div className="space-y-3">
              <div className="flex items-center justify-center gap-2 rounded-lg bg-green-500/10 px-4 py-3 text-sm text-green-400">
                <CheckCircle className="h-4 w-4" />
                Deposit confirmed!
              </div>
              {txHash && (
                <a
                  href={`https://sepolia.basescan.org/tx/${txHash}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center gap-1.5 text-xs text-blue-400 hover:underline"
                >
                  View on BaseScan <ExternalLink className="h-3 w-3" />
                </a>
              )}
              <Button className="w-full" onClick={() => setShowDialog(false)}>
                Done
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              <div>
                <Label className="text-xs text-muted-foreground">Amount (USDT)</Label>
                <Input
                  type="number"
                  min={0}
                  step={1}
                  value={amount}
                  onChange={(e) => setAmount(e.target.value)}
                  placeholder="e.g. 100"
                  className="mt-1"
                  disabled={isSending || isConfirming}
                />
              </div>
              <div className="flex gap-2">
                <Button
                  variant="outline"
                  className="flex-1"
                  onClick={() => setShowDialog(false)}
                  disabled={isSending || isConfirming}
                >
                  Cancel
                </Button>
                <Button
                  className="flex-1 gap-1.5"
                  onClick={handleDeposit}
                  disabled={!amount || isSending || isConfirming}
                >
                  {isSending ? (
                    <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Sending...</>
                  ) : isConfirming ? (
                    <><Loader2 className="h-3.5 w-3.5 animate-spin" /> Confirming...</>
                  ) : (
                    "Deposit"
                  )}
                </Button>
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>
    </>
  );
}

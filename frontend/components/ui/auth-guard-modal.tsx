"use client";

import { useAppKit } from "@reown/appkit/react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Wallet } from "lucide-react";

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function AuthGuardModal({ open, onOpenChange }: Props) {
  const { open: openWallet } = useAppKit();

  const handleConnect = () => {
    openWallet();
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader className="items-center text-center">
          <Wallet className="mb-2 h-10 w-10 text-primary" />
          <DialogTitle>Connect wallet required</DialogTitle>
          <DialogDescription>
            Connect your wallet to use this feature. Your wallet address is your identity on Miora.
          </DialogDescription>
        </DialogHeader>
        <div className="flex gap-2">
          <Button
            variant="outline"
            className="flex-1"
            onClick={() => onOpenChange(false)}
          >
            Cancel
          </Button>
          <Button className="flex-1 gap-2" onClick={handleConnect}>
            <Wallet className="h-4 w-4" />
            Connect Wallet
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

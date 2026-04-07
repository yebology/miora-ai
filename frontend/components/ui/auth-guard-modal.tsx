"use client";

import { useAuth } from "@/components/providers/auth-provider";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { LogIn, Loader2 } from "lucide-react";

type Props = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

export function AuthGuardModal({ open, onOpenChange }: Props) {
  const { signIn, loading } = useAuth();

  const handleSignIn = async () => {
    await signIn();
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-sm">
        <DialogHeader className="items-center text-center">
          <LogIn className="mb-2 h-10 w-10 text-primary" />
          <DialogTitle>Sign in required</DialogTitle>
          <DialogDescription>
            You need to sign in with Google to use this feature. It only takes a
            few seconds.
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
          <Button className="flex-1 gap-2" onClick={handleSignIn} disabled={loading}>
            {loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <LogIn className="h-4 w-4" />}
            Sign in
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}

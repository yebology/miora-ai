"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/providers/auth-provider";

// Redirect to analyze — auth is handled via wallet connect in navbar
export default function LoginPage() {
  const router = useRouter();
  const { isConnected } = useAuth();

  useEffect(() => {
    if (isConnected) {
      router.push("/analyze");
    }
  }, [isConnected, router]);

  return (
    <div className="flex flex-1 items-center justify-center py-24">
      <p className="text-sm text-muted-foreground">Connect your wallet to continue.</p>
    </div>
  );
}

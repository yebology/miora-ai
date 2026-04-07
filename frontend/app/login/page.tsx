"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/components/providers/auth-provider";

// Redirect to home — sign in is handled via navbar button
export default function LoginPage() {
  const router = useRouter();
  const { user, signIn } = useAuth();

  useEffect(() => {
    if (user) {
      router.push("/analyze");
    } else {
      signIn().then(() => router.push("/analyze"));
    }
  }, [user, signIn, router]);

  return (
    <div className="flex flex-1 items-center justify-center py-24">
      <div className="h-8 w-8 animate-spin rounded-full border-2 border-muted border-t-primary" />
    </div>
  );
}

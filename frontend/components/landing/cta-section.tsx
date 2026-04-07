"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useAnimateOnScroll } from "@/hooks/use-animate-on-scroll";

export function CtaSection() {
  const { ref, isVisible } = useAnimateOnScroll();

  return (
    <section className="relative overflow-hidden px-6 py-16">
      <div className="pointer-events-none absolute inset-0 bg-gradient-to-t from-primary/5 via-transparent to-transparent" />

      <div
        ref={ref}
        className={cn(
          "relative mx-auto flex max-w-md flex-col items-center gap-4 text-center transition-all duration-700 ease-out",
          isVisible ? "translate-y-0 opacity-100" : "translate-y-6 opacity-0"
        )}
      >
        <h2 className="text-2xl font-bold tracking-tight">Ready to start?</h2>
        <p className="text-sm text-muted-foreground">
          Analyze wallets and trade smarter — free to use.
        </p>
        <Link
          href="/analyze"
          className={cn(buttonVariants({ size: "lg" }), "gap-1.5")}
        >
          Launch App
          <ArrowRight className="h-4 w-4" />
        </Link>
      </div>
    </section>
  );
}

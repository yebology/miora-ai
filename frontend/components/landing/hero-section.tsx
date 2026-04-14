"use client";

import Link from "next/link";
import { ArrowRight } from "lucide-react";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import { useAnimateOnScroll } from "@/hooks/use-animate-on-scroll";
import { HeroBackground } from "@/components/landing/hero-background";

export function HeroSection() {
  const { ref, isVisible } = useAnimateOnScroll(0.1);

  return (
    <section className="relative overflow-hidden px-6 py-20 md:py-28">
      <HeroBackground />

      <div
        ref={ref}
        className={cn(
          "relative mx-auto flex max-w-5xl flex-col items-center gap-6 text-center transition-all duration-700 ease-out",
          isVisible ? "translate-y-0 opacity-100" : "translate-y-8 opacity-0"
        )}
      >
        <p className="rounded-full border bg-muted/50 px-4 py-1 text-xs text-muted-foreground backdrop-blur-sm">
          🧠 Trading Reputation Protocol on Base
        </p>

        <h1 className="max-w-3xl text-4xl font-bold tracking-tight sm:text-5xl md:text-6xl">
          Score any trader on Base.{" "}
          <span className="animate-shimmer bg-gradient-to-r from-purple-500 via-blue-500 to-cyan-500 bg-[length:200%_auto] bg-clip-text text-transparent">
            Let AI ride the winners for you.
          </span>
        </h1>

        <p className="max-w-lg text-muted-foreground">
          Analyze wallets, get AI-powered reputation scores, follow the best
          traders, and get notified when they move.
        </p>

        <div className="flex gap-3 pt-2">
          <Link
            href="/analyze"
            className={cn(buttonVariants({ size: "lg" }), "gap-1.5")}
          >
            Get Started
            <ArrowRight className="h-4 w-4" />
          </Link>
          <Link
            href="/swap"
            className={cn(buttonVariants({ variant: "outline", size: "lg" }))}
          >
            Swap Tokens
          </Link>
        </div>
      </div>
    </section>
  );
}

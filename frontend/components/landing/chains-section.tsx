"use client";

import Image from "next/image";
import { CHAINS } from "@/constants/landing";
import { cn } from "@/lib/utils";
import { useAnimateOnScroll } from "@/hooks/use-animate-on-scroll";

function ChainItem({
  name,
  logo,
  color,
}: {
  name: string;
  logo: string;
  color: string;
}) {
  return (
    <div
      className="group flex shrink-0 items-center gap-3 rounded-full border bg-card px-5 py-2.5 transition-all duration-300 hover:-translate-y-0.5 hover:scale-105"
      style={
        {
          "--chain-color": color,
        } as React.CSSProperties
      }
    >
      <Image
        src={logo}
        alt={name}
        width={24}
        height={24}
        className="h-6 w-6 transition-transform duration-300 group-hover:scale-110"
      />
      <span className="text-sm font-medium whitespace-nowrap transition-colors duration-300 group-hover:text-[var(--chain-color)]">
        {name}
      </span>
    </div>
  );
}

function MarqueeTrack() {
  return (
    <div className="flex shrink-0 gap-4">
      {CHAINS.map((chain) => (
        <ChainItem
          key={chain.name}
          name={chain.name}
          logo={chain.logo}
          color={chain.color}
        />
      ))}
    </div>
  );
}

export function ChainsSection() {
  const { ref, isVisible } = useAnimateOnScroll();

  return (
    <section className="bg-muted/30 py-16">
      <div
        ref={ref}
        className={cn(
          "transition-all duration-700 ease-out",
          isVisible ? "translate-y-0 opacity-100" : "translate-y-6 opacity-0"
        )}
      >
        <div className="mx-auto max-w-5xl px-6">
          <h2 className="text-center text-2xl font-bold tracking-tight">
            6 chains supported
          </h2>
          <p className="mt-2 mb-10 text-center text-sm text-muted-foreground">
            Wallet analysis + swap quotes on every chain.
          </p>
        </div>

        <div className="relative overflow-hidden py-4">
          <div className="pointer-events-none absolute inset-y-0 left-0 z-10 w-20 bg-gradient-to-r from-background to-transparent" />
          <div className="pointer-events-none absolute inset-y-0 right-0 z-10 w-20 bg-gradient-to-l from-background to-transparent" />

          <div className="flex w-fit animate-marquee gap-4 hover:[animation-play-state:paused]">
            <MarqueeTrack />
            <MarqueeTrack />
            <MarqueeTrack />
            <MarqueeTrack />
          </div>
        </div>
      </div>
    </section>
  );
}

"use client";

import { FEATURES } from "@/constants/landing";
import { cn } from "@/lib/utils";
import { useAnimateOnScroll } from "@/hooks/use-animate-on-scroll";

export function FeaturesSection() {
  const { ref, isVisible } = useAnimateOnScroll();

  return (
    <section className="bg-muted/30 px-6 py-16">
      <div ref={ref} className="mx-auto max-w-5xl">
        <h2 className="mb-10 text-center text-2xl font-bold tracking-tight">
          Everything you need
        </h2>

        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {FEATURES.map((feature, i) => (
            <div
              key={feature.title}
              className={cn(
                "group relative cursor-default overflow-hidden rounded-xl border bg-card p-5",
                "transition-all duration-300 ease-out",
                "hover:-translate-y-1 hover:scale-[1.02] hover:shadow-lg hover:shadow-primary/5",
                "hover:border-primary/30",
                isVisible
                  ? "translate-y-0 opacity-100"
                  : "translate-y-6 opacity-0"
              )}
              style={{ transitionDelay: isVisible ? "0ms" : `${i * 80}ms` }}
            >
              {/* Glow effect on hover */}
              <div className="pointer-events-none absolute inset-0 opacity-0 transition-opacity duration-300 group-hover:opacity-100">
                <div className="absolute -inset-1 bg-gradient-to-br from-purple-500/10 via-transparent to-cyan-500/10 blur-xl" />
              </div>

              <div className="relative">
                <feature.icon className="mb-3 h-5 w-5 text-muted-foreground transition-all duration-300 group-hover:scale-110 group-hover:text-primary" />
                <h3 className="text-sm font-medium transition-colors duration-300 group-hover:text-primary">
                  {feature.title}
                </h3>
                <p className="mt-1.5 text-sm leading-relaxed text-muted-foreground">
                  {feature.description}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

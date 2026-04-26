"use client";

import { STEPS } from "@/constants/landing";
import { cn } from "@/lib/utils";
import { useAnimateOnScroll } from "@/hooks/use-animate-on-scroll";

export function HowItWorksSection() {
  const { ref, isVisible } = useAnimateOnScroll();

  return (
    <section className="px-6 py-16">
      <div ref={ref} className="mx-auto max-w-5xl">
        <h2 className="mb-14 text-center text-2xl font-bold tracking-tight">
          How it works
        </h2>

        <div className="relative grid gap-8 md:grid-cols-3 md:gap-0">
          {/* Connector line (desktop only) */}
          <div className="pointer-events-none absolute left-0 right-0 top-10 hidden md:block">
            <div className="mx-auto h-px w-2/3 overflow-hidden bg-border">
              <div
                className={cn(
                  "h-full bg-gradient-to-r from-purple-500 via-blue-500 to-cyan-500 transition-all duration-1000 ease-out",
                  isVisible ? "w-full" : "w-0"
                )}
                style={{ transitionDelay: "400ms" }}
              />
            </div>
          </div>

          {STEPS.map((item, i) => (
            <div
              key={item.step}
              className={cn(
                "flex flex-col items-center text-center transition-all duration-500 ease-out",
                isVisible
                  ? "translate-y-0 opacity-100"
                  : "translate-y-8 opacity-0"
              )}
              style={{ transitionDelay: `${i * 200}ms` }}
            >
              {/* Icon circle with pulse */}
              <div className="relative mb-4">
                {/* Pulse ring */}
                <div
                  className={cn(
                    "absolute inset-0 rounded-full bg-primary/20",
                    isVisible && "animate-pulse-ring"
                  )}
                  style={{ animationDelay: `${i * 300 + 500}ms` }}
                />
                {/* Icon container */}
                <div className="relative flex h-20 w-20 items-center justify-center rounded-full border bg-card">
                  <item.icon className="h-7 w-7 text-muted-foreground" />
                </div>
              </div>

              <span className="text-xs font-medium text-muted-foreground/50">
                STEP {item.step}
              </span>
              <h3 className="mt-1 text-lg font-semibold">{item.title}</h3>
              <p className="mt-2 max-w-[220px] text-sm text-muted-foreground">
                {item.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

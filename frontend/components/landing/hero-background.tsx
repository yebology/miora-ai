"use client";

export function HeroBackground() {
  return (
    <div className="pointer-events-none absolute inset-0 overflow-hidden">
      {/* Animated gradient orbs — visible glow */}
      <div className="absolute -left-20 -top-20 h-[400px] w-[400px] animate-float rounded-full bg-purple-500/30 blur-[100px]" />
      <div className="absolute -right-20 top-10 h-[350px] w-[350px] animate-float-delayed rounded-full bg-blue-500/30 blur-[100px]" />
      <div className="absolute -bottom-10 left-1/3 h-[300px] w-[300px] animate-float-slow rounded-full bg-cyan-500/25 blur-[100px]" />

      {/* Subtle dot grid */}
      <div
        className="absolute inset-0 opacity-[0.04]"
        style={{
          backgroundImage:
            "radial-gradient(circle, currentColor 1px, transparent 1px)",
          backgroundSize: "40px 40px",
        }}
      />
    </div>
  );
}

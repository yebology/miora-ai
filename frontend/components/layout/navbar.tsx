"use client";

import { useState } from "react";
import Link from "next/link";
import { Brain, Menu, X, Wallet } from "lucide-react";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/layout/theme-toggle";
import { NAV_LINKS } from "@/constants/nav";
import { useAppKit, useAppKitAccount } from "@reown/appkit/react";

function shortenAddress(addr: string) {
  return `${addr.slice(0, 6)}...${addr.slice(-4)}`;
}

function ConnectWalletButton() {
  const { open } = useAppKit();
  const { address, isConnected } = useAppKitAccount();

  if (isConnected && address) {
    return (
      <button
        onClick={() => open()}
        className="flex h-8 items-center gap-1.5 rounded-lg border bg-card px-2.5 text-xs font-mono transition-colors hover:bg-muted/50"
      >
        <Wallet className="h-3 w-3 text-green-400" />
        {shortenAddress(address)}
      </button>
    );
  }

  return (
    <Button variant="outline" size="sm" className="gap-1.5" onClick={() => open()}>
      <Wallet className="h-3.5 w-3.5" />
      Connect Wallet
    </Button>
  );
}

export function Navbar() {
  const [open, setOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/80 backdrop-blur-sm">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-6">
        <Link href="/" className="flex items-center gap-2">
          <Brain className="h-5 w-5" />
          <span className="font-semibold">Miora AI</span>
        </Link>

        {/* Desktop nav */}
        <nav className="hidden items-center gap-6 text-sm md:flex">
          {NAV_LINKS.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              className="text-muted-foreground transition-colors hover:text-foreground"
            >
              {link.label}
            </Link>
          ))}
        </nav>

        <div className="hidden items-center gap-2 md:flex">
          <ThemeToggle />
          <ConnectWalletButton />
        </div>

        {/* Mobile toggle */}
        <div className="flex items-center gap-2 md:hidden">
          <ThemeToggle />
          <Button
            variant="ghost"
            size="icon"
            className="h-9 w-9"
            onClick={() => setOpen(!open)}
            aria-label="Toggle menu"
            aria-expanded={open}
          >
            {open ? <X className="h-4 w-4" /> : <Menu className="h-4 w-4" />}
          </Button>
        </div>
      </div>

      {/* Mobile menu */}
      {open && (
        <nav className="border-t px-6 py-4 md:hidden">
          <div className="flex flex-col gap-3">
            {NAV_LINKS.map((link) => (
              <Link
                key={link.href}
                href={link.href}
                onClick={() => setOpen(false)}
                className="text-sm text-muted-foreground transition-colors hover:text-foreground"
              >
                {link.label}
              </Link>
            ))}
            <div className="mt-2">
              <ConnectWalletButton />
            </div>
          </div>
        </nav>
      )}
    </header>
  );
}

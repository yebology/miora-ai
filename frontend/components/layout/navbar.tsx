"use client";

import { useState } from "react";
import Link from "next/link";
import { Brain, Menu, X, LogOut, Loader2 } from "lucide-react";
import { buttonVariants } from "@/components/ui/button";
import { Button } from "@/components/ui/button";
import { ThemeToggle } from "@/components/layout/theme-toggle";
import { NAV_LINKS } from "@/constants/nav";
import { useAuth } from "@/components/providers/auth-provider";
import { cn } from "@/lib/utils";

function UserMenu({
  user,
  onSignOut,
}: {
  user: { name: string; email: string; avatar: string };
  onSignOut: () => void;
}) {
  const [open, setOpen] = useState(false);

  return (
    <div className="relative">
      <button
        onClick={() => setOpen(!open)}
        className="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-xs font-medium text-primary-foreground transition-opacity hover:opacity-80"
      >
        {user.avatar}
      </button>

      {open && (
        <>
          <div
            className="fixed inset-0 z-40"
            onClick={() => setOpen(false)}
          />
          <div className="absolute right-0 top-full z-50 mt-2 w-56 rounded-lg border bg-popover p-1 shadow-lg">
            <div className="px-3 py-2">
              <p className="text-sm font-medium">{user.name}</p>
              <p className="text-xs text-muted-foreground">{user.email}</p>
            </div>
            <div className="my-1 h-px bg-border" />
            <button
              onClick={() => {
                onSignOut();
                setOpen(false);
              }}
              className="flex w-full items-center gap-2 rounded-md px-3 py-2 text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground"
            >
              <LogOut className="h-3.5 w-3.5" />
              Sign Out
            </button>
          </div>
        </>
      )}
    </div>
  );
}

export function Navbar() {
  const [open, setOpen] = useState(false);
  const { user, loading, signIn, signOut } = useAuth();

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
          {user ? (
            <UserMenu user={user} onSignOut={signOut} />
          ) : (
            <Button size="sm" onClick={signIn} disabled={loading}>
              {loading ? (
                <Loader2 className="h-3.5 w-3.5 animate-spin" />
              ) : (
                "Sign In"
              )}
            </Button>
          )}
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
            {user ? (
              <div className="mt-2 flex items-center justify-between rounded-lg bg-muted/50 px-3 py-2">
                <div className="flex items-center gap-2">
                  <div className="flex h-7 w-7 items-center justify-center rounded-full bg-primary text-xs font-medium text-primary-foreground">
                    {user.avatar}
                  </div>
                  <span className="text-sm">{user.name}</span>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-7 gap-1.5 text-xs text-muted-foreground"
                  onClick={() => {
                    signOut();
                    setOpen(false);
                  }}
                >
                  <LogOut className="h-3 w-3" />
                  Sign Out
                </Button>
              </div>
            ) : (
              <Button
                size="sm"
                className="mt-2 w-full"
                onClick={() => {
                  signIn();
                  setOpen(false);
                }}
                disabled={loading}
              >
                {loading ? "Signing in..." : "Sign In"}
              </Button>
            )}
          </div>
        </nav>
      )}
    </header>
  );
}

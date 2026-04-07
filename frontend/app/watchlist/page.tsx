"use client";

import { useState } from "react";
import type { WatchlistItem, Notification } from "@/types/watchlist";
import { DUMMY_WATCHLIST, DUMMY_NOTIFICATIONS } from "@/constants/dummy-watchlist";
import { WatchlistCard } from "@/components/watchlist/watchlist-card";
import { NotificationItem } from "@/components/watchlist/notification-item";
import { Eye, Bell } from "lucide-react";

export default function WatchlistPage() {
  const [watchlist, setWatchlist] = useState<WatchlistItem[]>(DUMMY_WATCHLIST);
  const [notifications] = useState<Notification[]>(DUMMY_NOTIFICATIONS);
  const [tab, setTab] = useState<"wallets" | "notifications">("wallets");

  const handleUnfollow = (address: string) => {
    setWatchlist((prev) => prev.filter((w) => w.wallet_address !== address));
  };

  const handleToggleNotify = (address: string) => {
    setWatchlist((prev) =>
      prev.map((w) =>
        w.wallet_address === address
          ? { ...w, email_notify: !w.email_notify }
          : w
      )
    );
  };

  const unreadCount = notifications.filter((n) => !n.read).length;

  return (
    <div className="px-6 py-10">
      <div className="mx-auto max-w-3xl">
        <div className="mb-8 text-center">
          <h1 className="text-2xl font-bold tracking-tight">Watchlist</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Wallets you follow and trade notifications.
          </p>
        </div>

        {/* Tabs */}
        <div className="mb-6 flex gap-1 rounded-lg bg-muted/50 p-1">
          <button
            onClick={() => setTab("wallets")}
            className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              tab === "wallets"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Eye className="h-4 w-4" />
            Wallets ({watchlist.length})
          </button>
          <button
            onClick={() => setTab("notifications")}
            className={`flex flex-1 items-center justify-center gap-2 rounded-md px-4 py-2 text-sm font-medium transition-colors ${
              tab === "notifications"
                ? "bg-background text-foreground shadow-sm"
                : "text-muted-foreground hover:text-foreground"
            }`}
          >
            <Bell className="h-4 w-4" />
            Notifications
            {unreadCount > 0 && (
              <span className="rounded-full bg-primary px-1.5 py-0.5 text-xs text-primary-foreground">
                {unreadCount}
              </span>
            )}
          </button>
        </div>

        {/* Content */}
        {tab === "wallets" && (
          <div className="space-y-3">
            {watchlist.length === 0 ? (
              <div className="py-16 text-center">
                <Eye className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
                <p className="text-sm text-muted-foreground">
                  No wallets followed yet. Analyze a wallet to start following.
                </p>
              </div>
            ) : (
              watchlist.map((item) => (
                <WatchlistCard
                  key={item.id}
                  item={item}
                  onUnfollow={handleUnfollow}
                  onToggleNotify={handleToggleNotify}
                />
              ))
            )}
          </div>
        )}

        {tab === "notifications" && (
          <div className="space-y-2">
            {notifications.length === 0 ? (
              <div className="py-16 text-center">
                <Bell className="mx-auto mb-3 h-8 w-8 text-muted-foreground/30" />
                <p className="text-sm text-muted-foreground">
                  No notifications yet. Follow a wallet to get trade alerts.
                </p>
              </div>
            ) : (
              notifications.map((n) => (
                <NotificationItem key={n.id} notification={n} />
              ))
            )}
          </div>
        )}
      </div>
    </div>
  );
}

export type NavItem = {
  label: string;
  href: string;
};

export const NAV_LINKS: NavItem[] = [
  { label: "Analyze", href: "/analyze" },
  { label: "Watchlist", href: "/watchlist" },
  { label: "Agent", href: "/agent" },
];

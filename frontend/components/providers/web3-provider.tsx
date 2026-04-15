"use client";

import { useRef } from "react";
import { createAppKit } from "@reown/appkit/react";
import { WagmiProvider } from "wagmi";
import { base } from "@reown/appkit/networks";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { WagmiAdapter } from "@reown/appkit-adapter-wagmi";
import type { ReactNode } from "react";

const queryClient = new QueryClient();

const projectId = process.env.NEXT_PUBLIC_REOWN_PROJECT_ID || "";

const metadata = {
  name: "Miora AI",
  description: "Trading reputation protocol + AI agent for Base",
  url: "https://miora.ai",
  icons: [],
};

const networks = [base] as const;

const wagmiAdapter = new WagmiAdapter({
  networks: [...networks],
  projectId,
  ssr: true,
});

let appKitInitialized = false;

function initAppKit() {
  if (appKitInitialized) return;
  appKitInitialized = true;

  createAppKit({
    adapters: [wagmiAdapter],
    networks: [...networks],
    projectId,
    metadata,
    features: {
      analytics: false,
    },
  });
}

export function Web3Provider({ children }: { children: ReactNode }) {
  const initialized = useRef(false);
  if (!initialized.current) {
    initAppKit();
    initialized.current = true;
  }

  return (
    <WagmiProvider config={wagmiAdapter.wagmiConfig}>
      <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
    </WagmiProvider>
  );
}

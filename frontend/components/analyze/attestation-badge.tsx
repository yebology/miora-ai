"use client";

import { ExternalLink, Shield } from "lucide-react";

type Props = {
  attestationUID?: string;
  explorerURL?: string;
};

export function AttestationBadge({ attestationUID, explorerURL }: Props) {
  if (!attestationUID) return null;

  return (
    <a
      href={explorerURL || `https://base-sepolia.easscan.org/attestation/view/${attestationUID}`}
      target="_blank"
      rel="noopener noreferrer"
      className="inline-flex items-center gap-1.5 rounded-full border border-blue-500/30 bg-blue-500/10 px-3 py-1 text-xs font-medium text-blue-400 transition-colors hover:bg-blue-500/20"
    >
      <Shield className="h-3 w-3" />
      Verified on Base
      <ExternalLink className="h-3 w-3" />
    </a>
  );
}

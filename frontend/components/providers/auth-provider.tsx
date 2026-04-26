"use client";

import { createContext, useContext, useState, useCallback, useEffect } from "react";
import { useAppKitAccount } from "@reown/appkit/react";

type User = {
  walletAddress: string;
  email?: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  walletAddress: string | undefined;
  isConnected: boolean;
};

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: false,
  walletAddress: undefined,
  isConnected: false,
});

export function useAuth() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { address, isConnected } = useAppKitAccount();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);

  // When wallet connects/disconnects, update user state
  useEffect(() => {
    if (isConnected && address) {
      setUser({ walletAddress: address.toLowerCase() });
    } else {
      setUser(null);
    }
  }, [isConnected, address]);

  return (
    <AuthContext.Provider value={{ user, loading, walletAddress: address, isConnected }}>
      {children}
    </AuthContext.Provider>
  );
}

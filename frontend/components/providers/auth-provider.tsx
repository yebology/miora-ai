"use client";

import { createContext, useContext, useState, useCallback } from "react";

type User = {
  name: string;
  email: string;
  avatar: string;
};

type AuthContextType = {
  user: User | null;
  loading: boolean;
  signIn: () => Promise<void>;
  signOut: () => void;
};

const AuthContext = createContext<AuthContextType>({
  user: null,
  loading: false,
  signIn: async () => {},
  signOut: () => {},
});

export function useAuth() {
  return useContext(AuthContext);
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(false);

  const signIn = useCallback(async () => {
    setLoading(true);
    try {
      // TODO: Replace with real Firebase Google sign-in
      // import { signInWithPopup, GoogleAuthProvider } from "firebase/auth";
      // const provider = new GoogleAuthProvider();
      // const result = await signInWithPopup(auth, provider);
      // const token = await result.user.getIdToken();
      // const res = await fetch(`${API_URL}/api/auth/me`, {
      //   headers: { Authorization: `Bearer ${token}` },
      // });
      // const json = await res.json();

      // Simulate sign-in
      await new Promise((r) => setTimeout(r, 1000));
      setUser({
        name: "Yobel Nathaniel",
        email: "yobel@example.com",
        avatar: "YN",
      });
    } finally {
      setLoading(false);
    }
  }, []);

  const signOut = useCallback(() => {
    // TODO: Firebase sign out
    // import { signOut as firebaseSignOut } from "firebase/auth";
    // await firebaseSignOut(auth);
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider value={{ user, loading, signIn, signOut }}>
      {children}
    </AuthContext.Provider>
  );
}

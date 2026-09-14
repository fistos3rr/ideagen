'use client';
import { createContext, useContext, useEffect, useState, ReactNode } from 'react';
import { authApi, meApi, setAccessToken, type components } from '@/lib/api';

type User = components['schemas']['dto.User'];

type Ctx = {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
};

const AuthCtx = createContext<Ctx | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const data = await meApi.get();
        setUser(data.user ?? null);
      } catch {
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const login = async (email: string, password: string) => {
    const { access_token } = await authApi.login({ email, password });
    setAccessToken(access_token ?? null);
    const me = await meApi.get();
    setUser(me.user ?? null);
  };

  const register = async (email: string, password: string) => {
    await authApi.register({ email, password });
  };

  const logout = async () => {
    try { await authApi.logout(); } finally {
      setAccessToken(null);
      setUser(null);
    }
  };

  return (
    <AuthCtx.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthCtx.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthCtx);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

'use client'

import { useAuth } from '@/lib/auth/useAuth';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

export default function ProtectedLayout({ children }: { children: React.ReactNode }) {
  const { user, loading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!loading && !user) router.replace('/login');
  }, [loading, user, router]);

  if (loading) {
    return <div style={{ padding: 24 }}>Загрузка…</div>;
  }
  if (!user) {
    return null; // редирект уже в процессе
  }
  return <>{children}</>;
}

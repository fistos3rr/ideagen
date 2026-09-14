// app/(protected)/me/page.tsx
'use client';

import { useAuth } from '@/lib/auth/useAuth';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

function formatDate(s?: string) {
  if (!s) return '—';
  try { return new Date(s).toLocaleString('ru-RU'); } catch { return s; }
}

export default function MePage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [busy, setBusy] = useState(false);

  async function onLogout() {
    setBusy(true);
    try {
      await logout();
    } finally {
      setBusy(false);
      router.replace('/login');
    }
  }

  if (!user) return null; // layout уже гарантирует, но на всякий

  return (
    <div className="me-page">
      <header className="me-header">
        <h1>Профиль</h1>
        <button onClick={onLogout} disabled={busy}>
          {busy ? 'Выходим…' : 'Выйти'}
        </button>
      </header>

      <dl className="me-card">
        <dt>ID</dt>          <dd>{user.id ?? '—'}</dd>
        <dt>Email</dt>       <dd>{user.email ?? '—'}</dd>
        <dt>Роль</dt>        <dd>{user.role ?? '—'}</dd>
        <dt>Создан</dt>      <dd>{formatDate(user.created_at)}</dd>
      </dl>
    </div>
  );
}

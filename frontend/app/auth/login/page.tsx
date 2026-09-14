'use client';

import { useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

export default function LoginPage() {
  const router = useRouter();
  const search = useSearchParams();
  const next = search.get('next') ?? '/me';

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    setLoading(false);

    if (!res.ok) {
      const payload = await res.json().catch(() => ({}));
      setError(
        typeof payload.error === 'string' ? payload.error : 'Неверный email или пароль',
      );
      return;
    }

    router.push(next);
    router.refresh(); // важно: обновить серверные компоненты
  }

  return (
    <form onSubmit={onSubmit} className="max-w-sm mx-auto space-y-4 p-8">
      <h1 className="text-2xl font-bold">Вход</h1>

      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
        className="w-full border rounded px-3 py-2"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Пароль"
        required
        className="w-full border rounded px-3 py-2"
      />

      {error && <p className="text-red-600 text-sm">{error}</p>}

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white rounded px-3 py-2 disabled:opacity-50"
      >
        {loading ? 'Входим…' : 'Войти'}
      </button>

      <p className="text-sm text-center">
        Нет аккаунта? <a href="/register" className="underline">Зарегистрироваться</a>
      </p>
    </form>
  );
}

'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function RegisterPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | Record<string, string> | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const res = await fetch('/api/auth/register', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ email, password }),
    });

    setLoading(false);

    if (!res.ok) {
      const payload = await res.json().catch(() => ({}));
      // бэк может вернуть { error: string } или { error: { email: '...' } }
      setError(payload.error ?? 'Ошибка регистрации');
      return;
    }

    // после регистрации обычно логинят сразу
    router.push('/login');
  }

  return (
    <form onSubmit={onSubmit} className="max-w-sm mx-auto space-y-4 p-8">
      <h1 className="text-2xl font-bold">Регистрация</h1>

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

      {error && (
        <div className="text-red-600 text-sm">
          {typeof error === 'string'
            ? error
            : Object.entries(error).map(([field, msg]) => (
                <p key={field}>{field}: {msg}</p>
              ))}
        </div>
      )}

      <button
        type="submit"
        disabled={loading}
        className="w-full bg-black text-white rounded px-3 py-2 disabled:opacity-50"
      >
        {loading ? 'Создаём…' : 'Зарегистрироваться'}
      </button>
    </form>
  );
}

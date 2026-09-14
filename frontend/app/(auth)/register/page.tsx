// app/(auth)/register/page.tsx
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/lib/auth/useAuth';
import { ApiError } from '@/lib/api';

export default function RegisterPage() {
  const { register, login } = useAuth();
  const router = useRouter();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [password2, setPassword2] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (password !== password2) {
      setError('Пароли не совпадают');
      return;
    }

    setLoading(true);
    try {
      await register(email, password);
      // сразу логинимся, чтобы не гонять пользователя на /login
      await login(email, password);
      router.replace('/me');
    } catch (err) {
      if (err instanceof ApiError) {
        // 422 от бэка — это map{field: message}
        const p = err.payload.error;
        if (typeof p === 'string') setError(p);
        else if (p && typeof p === 'object') {
          setError(Object.values(p).join(', '));
        } else setError('Ошибка регистрации');
      } else {
        setError('Не удалось связаться с сервером');
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth-page">
      <form className="auth-card" onSubmit={onSubmit}>
        <h1>Регистрация</h1>

        <label>
          Email
          <input
            type="email"
            autoComplete="email"
            required
            value={email}
            onChange={e => setEmail(e.target.value)}
          />
        </label>

        <label>
          Пароль
          <input
            type="password"
            autoComplete="new-password"
            required
            minLength={6}
            value={password}
            onChange={e => setPassword(e.target.value)}
          />
        </label>

        <label>
          Повторите пароль
          <input
            type="password"
            autoComplete="new-password"
            required
            value={password2}
            onChange={e => setPassword2(e.target.value)}
          />
        </label>

        {error && <p className="auth-error">{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? 'Создаём…' : 'Зарегистрироваться'}
        </button>

        <p className="auth-hint">
          Уже есть аккаунт? <Link href="/login">Войти</Link>
        </p>
      </form>
    </div>
  );
}

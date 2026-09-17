'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import type { LoginRequest } from '@/lib/api/types';
import type { FieldErrors } from '@/lib/api/error';

export function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<FieldErrors>({});
  const [pending, setPending] = useState(false);
  const router = useRouter()

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setValidationErrors({});
    setPending(true);

    try {
      const payload: LoginRequest = { email, password };
      const res = await fetch('/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      });

      const data = await res.json().catch(() => null);

      if (!res.ok) {
        if (res.status === 422 && data?.error) {
          const fields = data.error as FieldErrors;
          if (Object.keys(fields).length > 0) {
            setValidationErrors(fields);
            return;
          }
        }

        if (res.status === 401) {
          setError('Wrong email or password'); 
          return;
        }

        const message = typeof data?.error === 'string'
          ? data.error : `Error ${res.status}`;
        setError(message);
        return;
      } else {
        router.push('/me');
        router.refresh();
      } 
    } catch {
      setError('Try again later.');
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} noValidate>
      <h1>Login</h1>
      <label>
        Email
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
          autoComplete="email"
          aria-invalid={!!validationErrors.email}
          aria-describedby={validationErrors.email ? 'email-error' : undefined}
        />
      </label>
      {validationErrors.email && (
        <div id="email-error" role="alert" className="field-error">
          {validationErrors.email}
        </div>
      )}

      <label>
        Password
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          autoComplete="current-password"
          aria-invalid={!!validationErrors.password}
          aria-describedby={validationErrors.password ? 'password-error' : undefined}
        />
      </label>
      {validationErrors.password && (
        <div id="password-error" role="alert" className="field-error">
          {validationErrors.password}
        </div>
      )}

      {error && (
        <div role="alert" className="form-error">
          {error}
        </div>
      )}

      <button type="submit" disabled={pending}>
        {pending ? 'Entering...' : 'Enter'}
      </button>
    </form>
  );
}

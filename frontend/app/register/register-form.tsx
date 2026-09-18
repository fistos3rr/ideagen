'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import type { LoginRequest, UserCredentials } from '@/lib/api/types';
import type { FieldErrors } from '@/lib/api/error';

export function RegisterForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [repPassword, setRepPassword] = useState('');
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
      if (password !== repPassword) {
        setValidationErrors({ repPassword: 'Passwords must be same!'  });
        return
      } 

      const payload: UserCredentials = { email, password };
      const res = await fetch('/api/auth/register', {
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

        const message = typeof data?.error === 'string'
          ? data.error : `Error ${res.status}`;
        setError(message);
        return;
      } else {
        const loginReq: LoginRequest = { email, password };
        const loginRes = await fetch('/api/auth/login', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(loginReq),
        }).catch(() => null);
        if (!res.ok) {
          router.push('/login');
          router.refresh();
        } else {
          router.push('/');
          router.refresh();
        } 
      } 
    } catch {
      setError('Try again later.');
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} noValidate>
      <h1>Register</h1>
      <label>
        Email
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
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
        />
      </label>
      {validationErrors.password && (
        <div id="password-error" role="alert" className="field-error">
          {validationErrors.password}
        </div>
      )}

      <label>
        Repeat password
        <input
          type="password"
          value={repPassword}
          onChange={(e) => setRepPassword(e.target.value)}
          required
        />
      </label>
      {validationErrors.repPassword && (
        <div id="password-error" role="alert" className="field-error">
          {validationErrors.repPassword}
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

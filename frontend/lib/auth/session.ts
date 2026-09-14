import { cookies } from 'next/headers';
import { authApi, serviceApi } from '@/lib/api/endpoints';
import { ApiError } from '@/lib/api/client';

const ACCESS_COOKIE = 'access_token';
const ACCESS_MAX_AGE = 60 * 15;

export async function getAccessToken(): Promise<string | null> {
  const store = await cookies();
  return store.get(ACCESS_COOKIE)?.value ?? null;
}

export async function setAccessToken(token: string) {
  const store = await cookies();
  store.set(ACCESS_COOKIE, token, {
    httpOnly: true,
    secure: process.env.NODE_ENV === 'production',
    sameSite: 'lax',
    path: '/',
    maxAge: ACCESS_MAX_AGE,
  });
}

export async function clearAccessToken() {
  const store = await cookies();
  store.delete(ACCESS_COOKIE);
}

export async function getSession() {
  const token = await getAccessToken();
  if (!token) return null;

  try {
    const { user } = await serviceApi.me(token);
    return { user, token };
  } catch (e) {
    if (e instanceof ApiError && e.status === 401) {
      return null;
    }
    throw e;
  }
}

export async function requireSession() {
  const session = await getSession();
  if (!session) {
    const { redirect } = await import('next/navigation');
    redirect('/login');
  }
  return session;
}


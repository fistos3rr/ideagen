import 'server-only';
import { cookies } from 'next/headers';

const BACKEND = process.env.BACKEND_API_URL ?? 'http://localhost:4000';
const ACCESS_MAX_AGE = 60 * 15;

export async function apiRequest(path: string, init: RequestInit = {}) {
  const store = await cookies();
  let accessToken = store.get('access_token')?.value;
  const refreshToken = store.get('refresh_token')?.value;

  const doFetch = (token?: string) =>
    fetch(`${BACKEND}${path}`, {
      ...init,
      headers: {
        ...(init.headers || {}),
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
      },
    });

  let res = await doFetch(accessToken);

  if (res.status === 401 && refreshToken) {
    const refreshRes = await fetch(`${BACKEND}/auth/refresh`, {
      method: 'POST',
      headers: { Cookie: `refresh_token=${refreshToken}` },
    });

    if (refreshRes.ok) {
      const { access_token } = await refreshRes.json();
      accessToken = access_token;

      store.set('access_token', access_token, {
        httpOnly: true,
        secure: process.env.NODE_ENV === 'production',
        sameSite: 'lax',
        path: '/',
        maxAge: 60 * 15,
      });

      res = await doFetch(accessToken);
    } else {
      store.delete('access_token');
      store.delete('refresh_token');
    }
  }

  return res;
}

export function forwardSetCookie(from: Response, to: Response) {
  const cookies = from.headers.getSetCookie?.() ?? [];
  cookies.forEach((c) => to.headers.append('set-cookie', c));
}

import 'server-only';
import { cookies } from 'next/headers';
import { parseApiError } from './error';
import { BACKEND } from './auth';
import { decodeJwt } from 'jose';

// const ACCESS_MAX_AGE = 60 * 15;
// const REFRESH_MAX_AGE = 60 * 60 * 24;

type RequestOptions = {
  method?: 'POST' | 'GET' | 'PATCH' | 'DELETE';
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined>;
}

export function isTokenExpired(token: string): boolean {
  try {
    const payload = decodeJwt(token);
    if (!payload.exp) return true;

    const currentTime = Math.floor(Date.now() / 1000);
    return payload.exp <= currentTime + 5;
  } catch {
    return true;
  }
}

export async function isAuthorized(): boolean {
  const token = (await cookies()).get(ACCESS_COOKIE)?.value;
  return !isTokenExpired(token);
}

export async function apiRequest<T>(path: string, opts: RequestOptions = {}) {
  const store = await cookies();
  let accessToken = store.get('access_token')?.value;
  //const refreshToken = store.get('refresh_token')?.value;

  const url = new URL(`${BACKEND}${path}`);
  if (opts.query) {
    for (const [k, v] of Object.entries(opts.query)) {
      if (v !== undefined) {
        url.searchParams.set(k, String(v));
      }
    }
  }

  const doFetch = (token?: string) => {
    const headers: Record<string, string> = {};
    if (opts.body !== undefined) headers['Content-Type'] = 'application/json';
    if (token) headers['Authorization'] = `Bearer ${token}`

    return fetch(url.toString(), {
      method: opts.method ?? 'GET',
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
    });
  };

  let res = await doFetch(accessToken);

  // if (res.status === 401 && refreshToken) {
  //   const refreshRes = await fetch(`${BACKEND}/auth/refresh`, {
  //     method: 'POST',
  //     headers: { Cookie: `refresh_token=${refreshToken}`, Accept: 'application/json' },
  //     cache: 'no-store',
  //   });

  //   if (refreshRes.ok) {
  //     const { access_token } = await refreshRes.json();
  //     accessToken = access_token;

  //     store.set('access_token', access_token, {
  //       httpOnly: true,
  //       secure: process.env.NODE_ENV === 'production',
  //       sameSite: 'lax',
  //       path: '/',
  //       maxAge: ACCESS_MAX_AGE,
  //     });

  //     const newRefreshToken = extractRefreshToken(refreshRes);
  //     if (newRefreshToken) {
  //       store.set('refresh_token', newRefreshToken, {
  //         httpOnly: true,
  //         secure: process.env.NODE_ENV === 'production',
  //         sameSite: 'lax',
  //         path: '/',
  //         maxAge: REFRESH_MAX_AGE,
  //       });
  //     }

  //     res = await doFetch(accessToken);
  //   } else {
  //     store.delete('access_token');
  //     store.delete('refresh_token');
  //   }
  // }

  if (!res.ok) {
    throw await parseApiError(res);
  }

  if (res.status === 204) {
    return {} as T;
  }

  return res.json() as Promise<T>;
}

// function extractRefreshToken(res: Response): string | null {
//   const all = (res.headers as any).getSetCookie?.() ?? [];
//   for (const c of all) {
//     const m = /^refresh_token=([^;]+)/.exec(c);
//     if (m) return decodeURIComponent(m[1]);
//   }
//   return null
// }

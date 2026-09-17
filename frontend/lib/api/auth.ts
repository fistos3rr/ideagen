import type { NextResponse } from 'next/server';
import { cookies } from 'next/headers';
import { decodeJwt } from 'jose';
import { cache } from 'react';
import { serviceApi } from 'endpoints';

export const BACKEND = process.env.BACKEND_API_URL!;
export const IS_PROD = process.env.NODE_ENV === 'production';

export const ACCESS_TOKEN_TTL = 60 * 15;
export const REFRESH_TOKEN_TTL = 60 * 60 * 24;

export const ACCESS_COOKIE = 'access_token';
export const REFRESH_COOKIE = 'refresh_token';

const baseCookieOpts = {
  httpOnly: true,
  secure: IS_PROD,
  sameSite: 'lax' as const,
  path: '/',
};

export function setAccessCookie(res: NextResponse, token: string) {
  res.cookies.set(ACCESS_COOKIE, token, { ...baseCookieOpts, maxAge: ACCESS_TOKEN_TTL });
}

export function setRefreshCookie(res: NextResponse, token: string) {

  res.cookies.set(REFRESH_COOKIE, token, { 
    ...baseCookieOpts, 
    maxAge: REFRESH_TOKEN_TTL,
  });
}

export function clearAuthCookies(res: NextResponse) {
  res.cookies.set(ACCESS_COOKIE, '', {
    ...baseCookieOpts,
    path: '/',
    maxAge: 0
  });
  res.cookies.set(REFRESH_COOKIE, '', {
    ...baseCookieOpts,
    path: '/',
    maxAge: 0
  });
}

export function extractCookie(res: Response, name: string): string | null {
  const all = (res.headers as any).getSetCookie?.() ?? [];
  for (const c of all) {
    const m = new RegExp(`^${name}=([^;]+)`).exec(c);
    if (m) return decodeURIComponent(m[1]);
  }
  return null;
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


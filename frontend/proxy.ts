import { NextResponse, type NextRequest } from 'next/server';
import { LoginResponse } from '@/lib/api/types';
import { isTokenExpired } from '@/lib/api/client';
import { 
  BACKEND, setAccessCookie, 
  setRefreshCookie, extractCookie, clearAuthCookies,
  REFRESH_COOKIE, ACCESS_COOKIE
} from '@/lib/api/auth';

export const config = {
  matcher: ['/me'],
};

export async function proxy(req: NextRequest) {
  let accessToken = req.cookies.get(ACCESS_COOKIE)?.value;
  let refreshToken = req.cookies.get(REFRESH_COOKIE)?.value;

  if (accessToken && isTokenExpired(accessToken)) {
    accessToken = undefined;
  }

  if (!accessToken && !refreshToken) {
    return NextResponse.redirect(new URL('/login', req.url));
  }

  if (!accessToken && refreshToken) {
    try {
      const backRes = await fetch(`${BACKEND}/auth/refresh`, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          Cookie: `${REFRESH_COOKIE}=${refreshToken}`,
        },
        cache: 'no-store',
      });

      if (!backRes.ok) {
        const loginRes = NextResponse.redirect(new URL('/login', req.url));
        clearAuthCookies(loginRes);
        return loginRes;
      }
      

      const data = (await backRes.json()) as LoginResponse;
      if (!data.access_token) {
        return NextResponse.redirect(new URL('/login', req.url));
      }

      const refresh = extractCookie(backRes, REFRESH_COOKIE);
      if (!refresh) {
        return NextResponse.redirect(new URL('/login', req.url));
      }

      const res = NextResponse.next();
      setAccessCookie(res, data.access_token);
      setRefreshCookie(res, refresh);

      return res;
    } catch {
      return NextResponse.redirect(new URL('/login', req.url));  
    }
  }

  return NextResponse.next();
}


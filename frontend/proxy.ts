import { NextResponse, type NextRequest } from 'next/server';
import { LoginResponse } from '@/lib/api/types';
import { 
  BACKEND, setAccessCookie, isTokenExpired, 
  setRefreshCookie, extractCookie, clearAuthCookies,
  REFRESH_COOKIE, ACCESS_COOKIE
} from '@/lib/api/auth';

const PROTECTED_PATHS = ['/me']

function isProtected(pathname: string): boolean {
  return PROTECTED_PATHS.some(
    (p) => pathname === p || pathname.startsWith(p + '/')
  );
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|api).*)'],
};

export async function proxy(req: NextRequest) {
  const { pathname } = req.nextUrl; 

  let accessToken = req.cookies.get(ACCESS_COOKIE)?.value;
  let refreshToken = req.cookies.get(REFRESH_COOKIE)?.value;

  const authorized =
    Boolean(accessToken && !isTokenExpired(accessToken)) ||
    Boolean(refreshToken && !isTokenExpired(refreshToken));

  if (!authorized && isProtected(pathname)) {
    return NextResponse.redirect(new URL('/login', req.url));
  }

  if (!accessToken && refreshToken && isProtected(pathname)) {
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

  const result = NextResponse.next();
  result.headers.set('x-authorized', String(authorized));
  return result;
}


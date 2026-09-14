import { NextResponse, type NextRequest } from 'next/server';

const PROTECTED = ['/me'];
const AUTH_ONLY = ['/login', '/register'];

export function middleware(req: NextRequest) {
  const { pathname } = req.nextUrl;
  const hasRefresh = req.cookies.has('refresh_token');

  if (!hasRefresh && PROTECTED.some(p => pathname.startsWith(p))) {
    const url = req.nextUrl.clone();
    url.pathname = '/login';
    url.searchParams.set('next', pathname);
    return NextResponse.redirect(url);
  }

  if (hasRefresh && AUTH_ONLY.some(p => pathname.startsWith(p))) {
    const url = req.nextUrl.clone();
    url.pathname = '/me';
    url.search = '';
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/me/:path*', '/login', '/register'],
};

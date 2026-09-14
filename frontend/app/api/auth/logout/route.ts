import { NextResponse, type NextRequest } from 'next/server';
import { 
  BACKEND,
  clearAuthCookies,
  REFRESH_COOKIE,
} from '@/lib/api/auth';

export async function POST(req: NextRequest) {
  const currentRefreshToken = req.cookies.get(REFRESH_COOKIE)?.value;

  const res = NextResponse.json({
    ok: true,
  });

  clearAuthCookies(res);

  if (!currentRefreshToken) {
    return res
  }

  try {
    await fetch(`${BACKEND}/auth/logout`, {
      method: 'POST',
      headers: {
        Accept: 'application/json',
        Cookie: `${REFRESH_COOKIE}=${currentRefreshToken}`,
      },
      cache: 'no-store',
    });
  } catch {}

  return res;
}

import { NextResponse, type NextRequest } from 'next/server';
import { LoginResponse } from '@/lib/api/types';
import { parseApiError } from '@/lib/api/error';
import { 
  BACKEND, setAccessCookie, 
  setRefreshCookie, extractCookie,
  REFRESH_COOKIE,
} from '@/lib/api/auth';

// NOT USING ANYWHERE, refreshing through proxy.ts in root folder
export async function POST(req: NextRequest) {
  const currentRefreshToken = req.cookies.get(REFRESH_COOKIE)?.value;

  if (!currentRefreshToken) {
    return NextResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }

  const backRes = await fetch(`${BACKEND}/auth/refresh`, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      Cookie: `${REFRESH_COOKIE}=${currentRefreshToken}`,
    },
    cache: 'no-store',
  });

  if (!backRes.ok) {
    const err = await parseApiError(backRes);
    return NextResponse.json({ error: err.message }, { status: err.status });
  }

  const data = (await backRes.json()) as LoginResponse;
  if (!data.access_token) {
    return NextResponse.json({ error: 'Bad gateway' }, { status: 502 });
  }

  const refresh = extractCookie(backRes, REFRESH_COOKIE);
  if (!refresh) {
    return NextResponse.json({ error: 'Bad gateway' }, { status: 502 });
  }

  const res = NextResponse.json({
    ok: true,
  });
  setAccessCookie(res, data.access_token);
  setRefreshCookie(res, refresh);
  return res;
}

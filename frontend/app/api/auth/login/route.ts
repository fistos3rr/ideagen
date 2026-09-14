import { NextResponse, type NextRequest } from 'next/server';
import { LoginRequest, LoginResponse } from '@/lib/api/types';
import { parseApiError, FieldErrors } from '@/lib/api/error';
import { 
  BACKEND, setAccessCookie, 
  setRefreshCookie, extractCookie,
  REFRESH_COOKIE,
} from '@/lib/api/auth';

export async function POST(req: NextRequest) {
  const body = (await req.json().catch(() => null)) as LoginRequest | null;
  if (!body?.email || !body?.password) {
    return NextResponse.json(
      { error: 'email and password are required' },
      { status: 400 },
    );
  }

  const loginBody: LoginRequest = {
    email: body.email,
    password: body.password,
  }

  const backRes = await fetch(`${BACKEND}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
    body: JSON.stringify(loginBody),
    cache: 'no-store',
  });

  if (!backRes.ok) {
    const err = await parseApiError(backRes);
    if (err.status == 422) {
      const errors = err.payload.error as FieldErrors;
      return NextResponse.json({ error: errors } , { status: err.status })
    }
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

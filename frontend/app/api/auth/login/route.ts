import { NextResponse } from 'next/server';
import { authApi } from '@/lib/api/endpoints';
import { setAccessToken } from '@/lib/auth/session';
import { ApiError } from '@/lib/api/client';

export async function POST(req: Request) {
  const body = await req.json().catch(() => null);
  if (!body) {
    return NextResponse.json({ error: 'invalid body' }, { status: 400 });
  }

  try {
    const { access_token } = await authApi.login(body);

    // refresh_token уже улетел от бэка в Set-Cookie,
    // но нам важно, чтобы он осел на нашем домене.
    // Поэтому пробрасываем Set-Cookie от бэка в ответ клиенту.
    const res = NextResponse.json({ ok: true });

    // access_token — в наш httpOnly cookie
    res.cookies.set('access_token', access_token, {
      httpOnly: true,
      secure: process.env.NODE_ENV === 'production',
      sameSite: 'lax',
      path: '/',
      maxAge: 60 * 15,
    });

    return res;
  } catch (e) {
    if (e instanceof ApiError) {
      return NextResponse.json(e.payload, { status: e.status });
    }
    return NextResponse.json({ error: 'server error' }, { status: 500 });
  }
}

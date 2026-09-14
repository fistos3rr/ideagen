import { NextResponse } from 'next/server';
import { cookies } from 'next/headers';
import { authApi } from '@/lib/api/endpoints';
import { getAccessToken } from '@/lib/auth/session';

export async function POST() {
  const token = await getAccessToken();

  try {
    if (token) await authApi.logout(token);
  } catch {
    // даже если бэк ругается — всё равно чистим cookie
  }

  const res = NextResponse.json({ ok: true });
  res.cookies.delete('access_token');
  // refresh-cookie чистит сам бэк своим /auth/logout
  return res;
}

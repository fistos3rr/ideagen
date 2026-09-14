import { NextResponse } from 'next/server';
import { authApi } from '@/lib/api/endpoints';
import { ApiError } from '@/lib/api/client';

export async function POST(req: Request) {
  const body = await req.json().catch(() => null);
  if (!body) return NextResponse.json({ error: 'invalid body' }, { status: 400 });

  try {
    const result = await authApi.register(body);
    return NextResponse.json(result);
  } catch (e) {
    if (e instanceof ApiError) {
      return NextResponse.json(e.payload, { status: e.status });
    }
    return NextResponse.json({ error: 'server error' }, { status: 500 });
  }
}

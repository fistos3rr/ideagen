import { NextResponse, type NextRequest } from 'next/server';
import { UserCredentials, MessageResponse } from '@/lib/api/types';
import { parseApiError, FieldErrors } from '@/lib/api/error';
import { BACKEND } from '@/lib/api/auth';

export async function POST(req: NextRequest) {
  const body = (await req.json().catch(() => null)) as UserCredentials | null;
  if (body == null) {
    return NextResponse.json(
      { error: 'Request body cannot be empty' },
      { status: 400 },
    );
  }

  const creds: UserCredentials = {
    email: body.email,
    password: body.password,
  }

  const backRes = await fetch(`${BACKEND}/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Accept: 'application/json',
    },
    body: JSON.stringify(creds),
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

  const data = (await backRes.json()) as MessageResponse;

  const res = NextResponse.json({
    ok: true,
  });

  return res;
}

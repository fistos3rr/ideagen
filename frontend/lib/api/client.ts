import type { paths, components } from './schema';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? '/v1';

export class ApiError extends Error {
  constructor(
    public status: number,
    public payload: components['schemas']['dto.ErrorResponse'] | components['schemas']['dto.ValidationErrorResponse'],
  ) {
    super(typeof payload.error === 'string' ? payload.error : 'API error');
  }
}

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PATCH' | 'DELETE';
  body?: unknown;
  token?: string | null;
  query?: Record<string, string | number | boolean | undefined>;
};

export async function apiRequest<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const url = new URL(BASE_URL + path, window.location.origin);

  if (opts.query) {
    for (const [k, v] of Object.entries(opts.query)) {
      if (v !== undefined) url.searchParams.set(k, String(v));
    }
  }

  const headers: HeadersInit = {};
  if (opts.body) headers['Content-Type'] = 'application/json';
  if (opts.token) headers['Authorization'] = `Bearer ${opts.token}`;

  const res = await fetch(url.toString(), {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
    credentials: 'include',
  });

  if (!res.ok) {
    const payload = await res.json().catch(() => ({ error: 'unknown' }));
    throw new ApiError(res.status, payload);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

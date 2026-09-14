import type { paths, components } from './schema';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? '/v1';

type ApiErrorPayload =
  | components['schemas']['dto.ErrorResponse']
  | components['schemas']['dto.ValidationErrorResponse'];

export class ApiError extends Error {
  constructor(
    public status: number,
    public payload: ApiErrorPayload,
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
  const base = BASE_URL.startsWith('http')
    ? BASE_URL
    : window.location.origin + BASE_URL;
  const url = new URL(path.replace(/^\//, ''), base.endsWith('/') ? base : base + '/');

  if (opts.query) {
    for (const [k, v] of Object.entries(query)) {
      if (value !== undefined) url.searchParams.set(k, String(v));
    }
  }

  const headers: HeadersInit = {},
  if (opts.body) headers['Content-Type'] = 'application/json';
  if (opts.token) headers['Authorization'] = `Bearer ${opts.token}`;

  const res = await fetch(url.toString(), {
    method :opts.method ?? 'GET',
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

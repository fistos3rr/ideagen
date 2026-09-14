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

let accessToken: string | null = null;
let refreshPromise: Promise<string | null> | null = null;

export function setAccessToken(t: string | null) { accessToken = t; }
export function getAccessToken() { return accessToken; }

async function refreshAccess(): Promise<string | null> {
  if (refreshPromise) return refreshPromise;
  refreshPromise = (async () => {
    try {
      const res = await fetch(`${BASE_URL}/auth/refresh`, {
        method: 'POST',
        credentials: 'include',
      });
      if (!res.ok) return null;
      const data = (await res.json()) as components['schemas']['dto.LoginResponse'];
      accessToken = data.access_token ?? null;
      return accessToken;
    } catch {
      return null;
    } finally {
      refreshPromise = null;
    }
  })();
  return refreshPromise;
}

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PATCH' | 'DELETE';
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined>;
  _retry?: boolean;
};

export async function apiRequest<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const base = BASE_URL.startsWith('http')
    ? BASE_URL
    : window.location.origin + BASE_URL;
  const url = new URL(path.replace(/^\//, ''), base.endsWith('/') ? base : base + '/');

  if (opts.query) {
    for (const [k, v] of Object.entries(opts.query)) {
      if (v !== undefined) url.searchParams.set(k, String(v));
    }
  }

  const headers: HeadersInit = {};
  if (opts.body) headers['Content-Type'] = 'application/json';
  if (accessToken) headers['Authorization'] = `Bearer ${accessToken}`;

  const res = await fetch(url.toString(), {
    method: opts.method ?? 'GET',
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
    credentials: 'include',
  });

  if (res.status === 401 && !opts._retry) {
    const t = await refreshAccess();
    if (t) return apiRequest<T>(path, { ...opts, _retry: true });
  }

  if (!res.ok) {
    const payload = await res.json().catch(() => ({ error: 'unknown' }));
    throw new ApiError(res.status, payload);
  }

  if (res.status === 204) return undefined as T;
  return res.json();
}

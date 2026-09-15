import type { paths, components } from './schema';

const BACKEND_API_URL = process.env.BACKEND_API_URL ?? 'http://localhost:4000/v1';

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
export const setAccessToken = (t: string | null) => { accessToken = t; };
export const getAccessToken = () => accessToken;

let refreshPromise: Promise<string | null> | null = null;

function refreshAccessToken(): Promise<string | null> {
  if (refreshPromise) return refreshPromise;

  refreshPromise = (async () => {
    try {
      const res = await fetch(joinUrl(BACKEND_API_URL, '/auth/refresh'), {
        method: 'POST',
        credentials: 'include',
      });
      if (!res.ok) {
        accessToken = null;
        return null
      }
      const data = (await res.json()) as { accessToken: string };
      accessToken = data.accessToken;
      return accessToken;
    } catch {
      accessToken = null;
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

function joinUrl(base: string, path: string) {
  return new URL(base + path).toString();
}

function buildUrl(path: string, query?: RequestOptions['query']) {
  const url = new URL(joinUrl(BACKEND_API_URL, path));
  if (query) {
    for (const [k, v] of Object.entries(query)) {
      if (v !== undefined) url.searchParams.set(k, String(v));
    }
  }
  return url.toString();
}

async function safeJson(res: Response): Promise<ApiErrorPayload> {
  return res.json().catch(() => ({ error: 'unknown' })) as Promise<ApiErrorPayload>;
}

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PATCH' | 'DELETE';
  body?: unknown;
  query?: Record<string, string | number | boolean | undefined>;
  skipRefresh?: boolean;
};

export async function apiRequest<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const doFetch = (token: string | null) => {
    const headers: Record<string, string> = {};
    if (opts.body !== undefined) headers['Content-Type'] = 'application/json';
    if (token) headers['Authorization'] = `Bearer ${token}`;

    return fetch(buildUrl(path, opts.query), {
      method: opts.method ?? 'GET',
      headers,
      body: opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
      credentials: 'include',
    });
  };

  let res = await doFetch(getAccessToken());

  if (res.status === 401 && !opts.skipRefresh) {
    const newToken = await refreshAccessToken();
    if (!newToken) {
      throw new ApiError(401, await safeJson(res));
    }
    res = await doFetch(newToken);

    if (res.status === 401) {
      throw new ApiError(401, await safeJson(res));
    }
  }

  if (!res.ok) throw new ApiError(res.status, await safeJson(res));
  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}

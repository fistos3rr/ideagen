import { cookies } from 'next/headers';
import { apiRequest } from '@/lib/api/client';
import { authApi } from '@/lib/api/endpoints';
import { ApiError } from '@/lib/api/client';
import { getAccessToken, setAccessToken } from './session';

/**
 * Выполняет запрос к API с текущим access_token.
 * При 401 пытается обновить токен через /auth/refresh и повторить.
 */
export async function serverFetch<T>(
  fn: (token: string) => Promise<T>,
): Promise<T> {
  const token = await getAccessToken();
  if (!token) throw new ApiError(401, { error: 'unauthorized' });

  try {
    return await fn(token);
  } catch (e) {
    if (!(e instanceof ApiError) || e.status !== 401) throw e;

    // пробуем refresh
    const store = await cookies();
    const cookieHeader = store
      .getAll()
      .map((c) => `${c.name}=${c.value}`)
      .join('; ');

    try {
      const { access_token } = await authApi.refresh(cookieHeader);
      await setAccessToken(access_token);
      return await fn(access_token);
    } catch {
      // refresh тоже не сработал — кидаем 401 дальше
      throw new ApiError(401, { error: 'session expired' });
    }
  }
}

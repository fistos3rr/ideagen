import { apiRequest } from './client';
import type { components, paths } from './schema';

type LoginResponse = components['schemas']['dto.LoginResponse']

export const authApi = {
  login: (body: { email:string; password: string }) =>
    apiRequest<LoginResponse>('/auth/login', { method: 'POST', body })
}

import { apiRequest } from './client';
import type { components, paths } from './schema';

type Idea = components['schemas']['data.Idea'];
type IdeaListResponse = components['schemas']['dto.IdeaListResponse'];
type IdeaResponse = components['schemas']['dto.IdeaResponse'];
type Type = components['schemas']['data.Type'];
type TypeListResponse = components['schemas']['dto.TypeListResponse'];
type LoginResponse = components['schemas']['dto.LoginResponse'];
type UserResponse = components['schemas']['dto.UserResponse'];
type MessageResponse = components['schemas']['dto.MessageResponse'];
type HealthResponse = components['schemas']['dto.HealthResponse'];

type IdeasQuery = paths['/ideas']['get']['parameters']['query'];

export const healthApi = {
  health: async () => {
    const res = await apiRequest<HealthResponse>('/health', {
      method: 'GET',
    })
    return res;
  }
}

export const authApi = {
  login: async (body: { email: string; password: string }) => {
    const res = await apiRequest<LoginResponse>('/auth/login', {
      method: 'POST',
      body,
      skipRefresh: true,
    });
    setAccessToken(res.accessToken);
    return res;
  },

  refresh: async () => {
    const res = await apiRequest<LoginResponse>('/auth/refresh', {
      method: 'POST',
      skipRefresh: true,
    });
    setAccessToken(res.accessToken);
    return res;
  },

  register: (body: { email: string; password: string }) =>
    apiRequest<MessageResponse>('/auth/register', {
      method: 'POST',
      body,
      skipRefresh: true,
    }),

  logout: async () => {
    try {
      await apiRequest<MessageResponse>('/auth/logout', {
        method: 'POST',
        skipRefresh: true,
      });
    } finally {
      setAccessToken(null);
    }
  },
};

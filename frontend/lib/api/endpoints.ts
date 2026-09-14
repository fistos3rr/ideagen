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

type IdeasQuery = paths['/ideas']['get']['parameters']['query'];

export const authApi = {
  login: (body: { email: string; password: string }) =>
    apiRequest<LoginResponse>('/auth/login', { method: 'POST', body }),

  logout: (token: string) => 
    apiRequest<MessageResponse>('/auth/logout', { method: 'POST', token }),

  refresh: () =>
    apiRequest<LoginResponse>('/auth/refresh', { method: 'POST' }),

  register: (body: { email: string; password: string }) =>
    apiRequest<MessageResponse>('/auth/register', { method: 'POST', body }),
};

export const serviceApi = {
  me: (token: string) => apiRequest<UserResponse>('/service/me', { token }),
    
  myIdeas: (query?: IdeasQuery, token?: string) =>
    apiRequest<IdeaListResponse>('/service/ideas', { query, token }),

  myIdea: (id: number, token: string) =>
    apiRequest<IdeaResponse>(`/service/ideas/${id}`, { token }),

  generateIdea: (token: string) =>
  apiRequest<components['schemas']['dto.BufferIdeaResponse']>(
    '/service/idea/generate',
    { method: 'POST', token },
  ),

  bufferIdeas: (token: string) =>
  apiRequest<components['schemas']['dto.BufferIdeaListResponse']>(
    '/service/idea/buffer',
    { token },
  ),

  chooseBufferIdea: (buffer_idea_id: string, token: string) =>
  apiRequest<IdeaResponse>('/service/idea/buffer', {
    method: 'POST',
    body: { buffer_idea_id },
    token,
  }),
};

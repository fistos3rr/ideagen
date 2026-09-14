import { apiRequest } from './client';
import type { components, paths } from './schema';

type S = components['schemas'];

export const authApi = {
  login: (body: S['dto.LoginRequest']) =>
    apiRequest<S['dto.LoginResponse']>('/auth/login', { method: 'POST', body }),

  register: (body: S['dto.UserCredentials']) =>
    apiRequest<S['dto.MessageResponse']>('/auth/register', { method: 'POST', body }),

  logout: () =>
    apiRequest<S['dto.MessageResponse']>('/auth/logout', { method: 'POST' }),

  refresh: () =>
    apiRequest<S['dto.LoginResponse']>('/auth/refresh', { method: 'POST' }),
};

export const meApi = {
  get: () => apiRequest<S['dto.UserResponse']>('/service/me'),
};

// Query-типы достаём прямо из paths — они уже описаны в схеме
type IdeasQuery = paths['/ideas']['get']['parameters']['query'];

export const ideasApi = {
  list: (query?: IdeasQuery) =>
    apiRequest<S['dto.IdeaListResponse']>('/ideas', { query }),

  get: (id: number) =>
    apiRequest<S['dto.IdeaResponse']>(`/ideas/${id}`),

  create: (body: S['dto.IdeaRequest']) =>
    apiRequest<S['dto.IdeaResponse']>('/ideas', { method: 'POST', body }),

  update: (id: number, body: S['dto.IdeaUpdateRequest']) =>
    apiRequest<S['dto.IdeaResponse']>(`/ideas/${id}`, { method: 'PATCH', body }),

  remove: (id: number) =>
    apiRequest<S['dto.MessageResponse']>(`/ideas/${id}`, { method: 'DELETE' }),
};

export const typesApi = {
  list: (query?: paths['/types']['get']['parameters']['query']) =>
    apiRequest<S['dto.TypeListResponse']>('/types', { query }),

  get: (id: number) =>
    apiRequest<S['dto.TypeResponse']>(`/types/${id}`),

  create: (body: S['dto.TypeRequest']) =>
    apiRequest<S['dto.TypeResponse']>('/types', { method: 'POST', body }),

  update: (id: number, body: S['dto.TypeUpdateRequest']) =>
    apiRequest<S['dto.TypeResponse']>(`/types/${id}`, { method: 'PATCH', body }),

  remove: (id: number) =>
    apiRequest<S['dto.MessageResponse']>(`/types/${id}`, { method: 'DELETE' }),
};

export const bufferApi = {
  list: () => apiRequest<S['dto.BufferIdeaListResponse']>('/service/idea/buffer'),

  get: (uuid: string) =>
    apiRequest<S['dto.BufferIdeaResponse']>(`/service/idea/buffer/${uuid}`),

  generate: () =>
    apiRequest<S['dto.BufferIdeaResponse']>('/service/idea/generate', { method: 'POST' }),

  choose: (buffer_idea_id: string) =>
    apiRequest<S['dto.IdeaResponse']>('/service/idea/buffer', {
      method: 'POST',
      body: { buffer_idea_id } satisfies S['dto.BufferIdeaRequest'],
    }),
};

export const serviceApi = {
  listIdeas: (query?: paths['/service/ideas']['get']['parameters']['query']) =>
    apiRequest<S['dto.IdeaListResponse']>('/service/ideas', { query }),

  getIdea: (id: number) =>
    apiRequest<S['dto.IdeaResponse']>(`/service/ideas/${id}`),

  removeIdea: (id: number) =>
    apiRequest<S['dto.MessageResponse']>(`/service/ideas/${id}`, { method: 'DELETE' }),

  bindUserIdea: (body: S['dto.UserIdeaRequest']) =>
    apiRequest<S['dto.MessageResponse']>('/service/useridea', { method: 'POST', body }),
};

export const aiApi = {
  ask: (message: string) =>
    apiRequest<S['dto.AskResponse']>('/ask', {
      method: 'POST',
      body: { message } satisfies S['dto.AskRequest'],
    }),
};

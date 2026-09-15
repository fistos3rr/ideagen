import { apiRequest } from "./client";
import type { components, paths } from "./schema.d.ts";

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

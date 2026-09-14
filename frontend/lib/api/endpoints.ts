import { apiRequest } from "./client";
import type { HealthResponse, UserResponse } from "./types";

export const healthApi = {
  health: async () => {
    const res = await apiRequest<HealthResponse>('/health', {
      method: 'GET',
    })
    return res;
  }
}

export const serviceApi = {
  me: async () => {
    const res = await apiRequest<UserResponse>('/service/me', {
      method: 'GET',
    })
    return res;
  }
}

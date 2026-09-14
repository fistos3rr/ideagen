import type { components, paths } from "./schema.d.ts";

export type IdeasQuery = paths['/ideas']['get']['parameters']['query'];

type S = components['schemas'];

export type HealthResponse = S['dto.HealthResponse'];
export type LoginRequest = S['dto.LoginRequest'];
export type LoginResponse = S['dto.LoginResponse'];
export type ErrorResponse = S['dto.ErrorResponse'];
export type ValidationErrorResponse = S['dto.ValidationErrorResponse'];
export type MessageResponse = S['dto.MessageResponse'];
export type User = S['dto.User'];
export type UserResponse = S['dto.UserResponse'];

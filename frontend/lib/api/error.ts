import type { ErrorResponse, ValidationErrorResponse } from "./types";

type ApiErrorPayload = ErrorResponse | ValidationErrorResponse;
export type FieldErrors = Record<string, string>;

export class ApiError extends Error {
  status: number;
  payload: ApiErrorPayload;

  constructor(opts: {
    status: number;
    payload: ApiErrorPayload;
    message?: string;
  }) {
    const { status, payload, message } = opts;
    let msg = message ?? defaultMessage(status);
    super(msg);
    this.name = 'ApiError';
    this.status = status;
    this.payload = payload;
  }
}

export async function parseApiError(res: Response): Promise<ApiError> {
  let body: any = null;
  const text = await res.text().catch(() => '');
  if (text) {
    try { body = JSON.parse(text); } catch { 
      return new ApiError({
        status: 500,
        payload: { error: "JSON Parse error" },
      }) 
    }
  }

  // 422: validation error
  if (res.status === 422 && body?.error && typeof body.error === 'object') {
    const fields = body.error as FieldErrors;
    const first = Object.values(fields)[0];
    return new ApiError({
      status: 422,
      message: first,
      payload: body,
    });
  }

  return new ApiError({ status: res.status, payload: body ?? { error: defaultMessage(res.status) } });
}

function defaultMessage(status: number): string {
  switch (status) {
    case 400: return 'Bad request';
    case 401: return 'Unauthorized';
    case 403: return 'Forbidden';
    case 404: return 'Not found';
    case 422: return 'Validation error';
    case 500: return 'Internal server error';
    default:  return `Error ${status}`;
  }
}

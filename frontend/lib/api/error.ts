export type FieldErrors = Record<string, string>;

export class ApiError extends Error {
  status: number;
  fields?: FieldErrors;
  raw?: unknown;
  
  constructor(opts: {
    status: number;
    message: string;
    fields?: FieldErrors;
    raw?: unknown;
  }) {
    super(opts.message);
    this.name = 'ApiError';
    this.status = opts.status;
    this.fields = opts.fields;
    this.raw = opts.raw;
  }

  get isUnauthorized() { return this.status === 401; }
  get isForbidden()    { return this.status === 403; }
  get isNotFound()     { return this.status === 404; }
  get isValidation()   { return this.status === 422; }
  get isServer()       { return this.status >= 500; }
}

export async function parseApiError(res: Response): Promise<ApiError> {
  let body: any = null;
  const text = await res.text().catch(() => '');
  if (text) {
    try { body = JSON.parse(text); } catch { /* не JSON */ }
  }

  // 422: validation error
  if (res.status === 422 && body?.error && typeof body.error === 'object') {
    const fields = body.error as FieldErrors;
    const first = Object.values(fields)[0] ?? 'Validation error';
    return new ApiError({
      status: 422,
      message: first,
      fields,
      raw: body,
    });
  }

  const message =
    (typeof body?.error === 'string' && body.error) ||
    defaultMessage(res.status);

  return new ApiError({ status: res.status, message, raw: body });
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

import { SERVICE_TOKEN } from '@shared/config/service';
import { signInternalRequest } from '@shared/crypto/hmac';
import createClient, { type Middleware } from 'openapi-fetch';
import type { paths } from './schema';
import type { ErrorResponse } from './types';

export class ApiError extends Error {
  readonly status: number;
  readonly code: string | undefined;
  readonly requestId: string | undefined;

  constructor(status: number, message: string, code?: string, requestId?: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.requestId = requestId;
  }
}

type TokenGetter = () => string | null;

let tokenGetter: TokenGetter = () => null;
let onUnauthorized: (() => void) | null = null;

export const setTokenGetter = (fn: TokenGetter): void => {
  tokenGetter = fn;
};

export const setOnUnauthorized = (fn: () => void): void => {
  onUnauthorized = fn;
};

let classifiedConfig: Record<string, unknown> = {};

export const setClassifiedConfig = (cfg: Record<string, unknown>): void => {
  classifiedConfig = cfg;
};

const toHeaderKey = (k: string): string => `X-Config-${k.replace(/[^a-zA-Z0-9-]/g, '-')}`;

const toHeaderValue = (v: unknown): string => {
  if (v === null || v === undefined) return '';
  if (typeof v === 'string') return v;
  if (typeof v === 'number' || typeof v === 'boolean') return String(v);
  try {
    return JSON.stringify(v);
  } catch {
    return '';
  }
};

const authMiddleware: Middleware = {
  onRequest: async ({ request }) => {
    const token = tokenGetter();
    if (token) {
      request.headers.set('Authorization', `Bearer ${token}`);
    }
    request.headers.set('Accept', 'application/json');

    const url = new URL(request.url);
    if (url.pathname.startsWith('/api/archive')) {
      for (const key in classifiedConfig) {
        const value = (classifiedConfig as Record<string, unknown>)[key];
        if (typeof value === 'function' || isObjectLike(value)) continue;
        request.headers.set(toHeaderKey(key), toHeaderValue(value));
      }
    }

    if (url.pathname.startsWith('/api/internal')) {
      if (SERVICE_TOKEN) {
        request.headers.set('X-Service-Token', SERVICE_TOKEN);
      }
      const signature = await signInternalRequest(
        `${request.method.toUpperCase()} ${url.pathname}`,
      );
      request.headers.set('X-Internal-Signature', signature);
    }

    return request;
  },
  onResponse: ({ response }) => {
    if (response.status === 401 && onUnauthorized) {
      onUnauthorized();
    }
    return response;
  },
};

const isObjectLike = (v: unknown): boolean => typeof v === 'object' && v !== null;

export const api = createClient<paths>({ baseUrl: '' });
api.use(authMiddleware);

type FetchResult<T> = {
  data?: T;
  error?: unknown;
  response: Response;
};

export const unwrap = async <T>(promise: Promise<FetchResult<T>>): Promise<T> => {
  const result = await promise;
  if (result.error !== undefined || result.data === undefined) {
    const err = result.error as Partial<ErrorResponse> | undefined;
    const status = result.response.status;
    const message = err?.error?.message ?? `Request failed (${status})`;
    throw new ApiError(status, message, err?.error?.code, err?.request_id);
  }
  return result.data;
};

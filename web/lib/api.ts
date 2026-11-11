import type { ApiError } from "./types";

// Server components call the Go API directly; the browser goes through the
// same-origin rewrite in next.config.ts (/api/* → API_URL/api/*).
const SERVER_BASE = process.env.API_URL ?? "http://localhost:8080";

export class ApiRequestError extends Error {
  constructor(
    public status: number,
    public code: string,
    message: string,
  ) {
    super(message);
  }
}

async function parse<T>(res: Response): Promise<T> {
  if (res.ok) return (await res.json()) as T;
  let body: ApiError = { error: "unknown", message: res.statusText };
  try {
    body = (await res.json()) as ApiError;
  } catch {
    /* non-JSON error body */
  }
  throw new ApiRequestError(res.status, body.error, body.message);
}

/**
 * Server-side GET, cached by Next.js for an hour (the API also caches).
 * Times out so a slow or sleeping API can never hang a render or a build.
 */
export async function serverGet<T>(path: string, revalidate = 3600, timeoutMs = 10000): Promise<T> {
  const res = await fetch(`${SERVER_BASE}${path}`, {
    next: { revalidate },
    signal: AbortSignal.timeout(timeoutMs),
  });
  return parse<T>(res);
}

/** Browser-side request through the same-origin rewrite. */
export async function clientFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    ...init,
    headers: { "Content-Type": "application/json", ...(init?.headers ?? {}) },
  });
  return parse<T>(res);
}

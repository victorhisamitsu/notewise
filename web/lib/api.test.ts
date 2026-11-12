import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiRequestError, clientFetch } from "./api";

afterEach(() => vi.unstubAllGlobals());

describe("clientFetch", () => {
  it("returns the JSON body on success", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => new Response(JSON.stringify({ ok: true }), { status: 200 })));
    await expect(clientFetch("/api/v1/families")).resolves.toEqual({ ok: true });
  });

  it("throws a typed error with the API's code and message", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => new Response(JSON.stringify({ error: "not_found", message: "note not found" }), { status: 404 })),
    );
    const err = await clientFetch("/api/v1/notes/x").catch((e) => e);
    expect(err).toBeInstanceOf(ApiRequestError);
    expect(err).toMatchObject({ status: 404, code: "not_found", message: "note not found" });
  });

  it("copes with a non-JSON error body", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => new Response("<html>bad gateway</html>", { status: 502, statusText: "Bad Gateway" })));
    const err = await clientFetch("/api/v1/families").catch((e) => e);
    expect(err).toMatchObject({ status: 502, code: "internal_error", message: "Bad Gateway" });
  });
});

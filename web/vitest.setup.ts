import "@testing-library/jest-dom/vitest";
import { cleanup } from "@testing-library/react";
import { afterEach, vi } from "vitest";
import { createElement } from "react";

afterEach(() => cleanup());

// next/image needs the Next.js runtime; in unit tests a plain <img> is enough.
vi.mock("next/image", () => ({
  default: ({ src, alt, fill: _fill, priority: _priority, ...rest }: Record<string, unknown>) =>
    createElement("img", { src, alt, ...rest }),
}));

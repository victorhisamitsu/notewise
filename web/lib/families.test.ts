import { describe, expect, it } from "vitest";
import { FALLBACK_COLOR, familyColor, shade } from "./families";

describe("shade", () => {
  it("mixes towards white and black", () => {
    expect(shade("#000000", 0.5)).toBe("#808080");
    expect(shade("#ffffff", -1)).toBe("#000000");
    expect(shade("#336699", 0)).toBe("#336699");
  });
});

describe("familyColor", () => {
  it("falls back for unknown families", () => {
    expect(familyColor("sunlit-citrus")).not.toBe(FALLBACK_COLOR);
    expect(familyColor("not-a-family")).toBe(FALLBACK_COLOR);
    expect(familyColor(undefined)).toBe(FALLBACK_COLOR);
  });
});

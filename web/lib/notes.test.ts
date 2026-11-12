import { existsSync } from "node:fs";
import { describe, expect, it } from "vitest";
import photos from "./note-photos.json";
import { notePhoto, photoKey, slugify } from "./notes";

describe("slugify", () => {
  it.each([
    ["Tonka Bean", "tonka-bean"],
    ["Ô de Lancôme", "o-de-lancome"],
    ["J'adore", "jadore"],
    ["  Pink   Pepper ", "pink-pepper"],
  ])("%s → %s", (input, want) => expect(slugify(input)).toBe(want));
});

describe("photoKey", () => {
  it.each([
    ["Vanilla", "vanilla"],
    ["Calabrian Bergamot", "bergamot"],
    ["Damask Rose", "rose"],
    ["Vanilla Absolute", "vanilla"],
    ["Cedarwood", "cedar"],
    ["Orange Blossom", "neroli"],
    ["Tonka", "tonka-bean"],
    ["Agarwood", "oud"],
  ])("%s → %s", (name, key) => expect(photoKey(name)).toBe(key));

  it("returns null for notes without a photo", () => {
    expect(photoKey("Guaiac Wood")).toBeNull();
    expect(notePhoto("Ambroxan")).toBeNull();
  });
});

describe("note photos", () => {
  it("every credited photo exists and every credit has an author and a license", () => {
    for (const [key, p] of Object.entries(photos)) {
      expect(existsSync(`public/notes/${key}.webp`), key).toBe(true);
      expect(p.author, key).not.toBe("");
      expect(p.license, key).not.toBe("");
      expect(p.source, key).toMatch(/^https:\/\/commons\.wikimedia\.org\//);
    }
  });
});

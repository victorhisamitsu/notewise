import photos from "./note-photos.json";

export type NotePhoto = {
  title: string;
  author: string;
  license: string;
  licenseUrl: string;
  source: string;
};

const PHOTOS = photos as Record<string, NotePhoto>;

/** "Tonka Bean" → "tonka-bean" (same rule as the Go API's Slugify). */
export function slugify(s: string): string {
  return s
    .toLowerCase()
    .normalize("NFD")
    .replace(/[̀-ͯ]/g, "")
    .replace(/['’]/g, "")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

// Names perfumers use for the same ingredient.
const ALIASES: Record<string, string> = {
  cedarwood: "cedar",
  "virginia-cedar": "cedar",
  "atlas-cedar": "cedar",
  agarwood: "oud",
  oudh: "oud",
  frankincense: "incense",
  olibanum: "incense",
  tonka: "tonka-bean",
  "jasmine-sambac": "jasmine",
  "orange-blossom": "neroli",
  "bitter-orange": "orange",
  "sweet-orange": "orange",
  "blood-orange": "orange",
  cacao: "cocoa",
  chocolate: "cocoa",
  labdanum: "amber",
  "black-tea": "tea",
  "green-tea": "tea",
  tangerine: "mandarin",
  "pink-peppercorn": "pink-pepper",
  peppermint: "mint",
  spearmint: "mint",
  "lily-of-the-valley": "lily-of-the-valley",
  muguet: "lily-of-the-valley",
};

/** Finds the photo key for a note name, e.g. "Calabrian Bergamot" → "bergamot". */
export function photoKey(name: string): string | null {
  const slug = slugify(name);
  if (PHOTOS[slug]) return slug;
  if (ALIASES[slug] && PHOTOS[ALIASES[slug]]) return ALIASES[slug];
  const words = slug.split("-");
  // "vanilla-absolute" → vanilla, "damask-rose" → rose: try the last word, then the first.
  for (const w of [words.at(-1), words[0]]) {
    if (!w) continue;
    if (PHOTOS[w]) return w;
    if (ALIASES[w] && PHOTOS[ALIASES[w]]) return ALIASES[w];
  }
  return null;
}

export function notePhoto(name: string): { src: string; credit: NotePhoto } | null {
  const key = photoKey(name);
  return key ? { src: `/notes/${key}.webp`, credit: PHOTOS[key] } : null;
}

export function allNotePhotos(): (NotePhoto & { key: string })[] {
  return Object.entries(PHOTOS).map(([key, p]) => ({ key, ...p }));
}

/** Ingredients shown on the home page and in the finder. */
export const FEATURED_NOTES = [
  "Vanilla", "Bergamot", "Rose", "Oud", "Sandalwood", "Jasmine",
  "Tonka Bean", "Lavender", "Iris", "Pink Pepper", "Saffron", "Coffee",
];

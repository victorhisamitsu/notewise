// Contracts returned by the Notewise Go API (api/internal/.../dto).

export type Family = {
  id: string;
  name: string;
  definition?: string;
  count?: number;
  subfamilies?: string[];
};

export type Summary = {
  id: string;
  name: string;
  house: string;
  family: Family;
};

export type Pyramid = { top: string[]; heart: string[]; base: string[] };

export type Perfume = Summary & {
  year?: number;
  perfumer?: string;
  concentration?: string;
  subfamily?: string;
  pyramid?: Pyramid;
  notes: string[];
  character?: string;
  characterWords: string[];
  seasons: string[];
  hours: string[];
  signatureNotes: string[];
  sourceUrl?: string;
};

export type Neighbour = Summary & { sharedNotes: string[]; difference?: string };
export type NoteMatch = Summary & { tiers: Record<string, string> };
export type Pick = Summary & { role?: string; why: string };

export type SearchResponse = { query: string; total: number; results: Summary[]; attribution: string };
export type PerfumeResponse = { perfume: Perfume; neighbours: Neighbour[]; attribution: string };
export type NoteMatchResponse = { notes: string[]; tier: string; total: number; results: NoteMatch[]; attribution: string };
export type FamiliesResponse = { families: Family[]; attribution: string };
export type TodayResponse = {
  occasion: string;
  headline: string;
  reason: string;
  mood: string;
  picks: Pick[];
  attribution: string;
};
export type Count = { label: string; count: number };
export type ProfileResponse = {
  size: number;
  families: Count[];
  topNotes: Count[];
  seasons: Count[];
  hours: Count[];
  characterWords: Count[];
  headline: string;
  gaps: string[];
  suggestions: Pick[];
  missing: string[];
  attribution: string;
};
export type ApiError = { error: string; message: string };

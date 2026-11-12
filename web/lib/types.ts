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
export type Note = {
  id: string;
  name: string;
  kind?: string;
  smellsLike?: string;
  origin?: string;
  category?: string;
  family: Family;
  pairsWith: string[];
  aliases: string[];
  showcasedIn: NoteMatch[];
};
export type NoteResponse = { note: Note; perfumes: NoteMatch[]; total: number; attribution: string };
export type ApiError = { error: string; message: string };

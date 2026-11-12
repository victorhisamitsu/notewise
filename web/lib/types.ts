// API contracts, generated from api/openapi.yaml (see lib/api-schema.ts).
// Never edit these shapes by hand: change the OpenAPI file and run
// `npm run gen:types`. CI fails if the generated file is out of date.
import type { components } from "./api-schema";

type S = components["schemas"];

export type Family = S["Family"];
export type Summary = S["Summary"];
export type Pyramid = S["Pyramid"];
export type Perfume = S["Perfume"];
export type Neighbour = S["Neighbour"];
export type NoteMatch = S["NoteMatch"];
export type Pick = S["Pick"];
export type Note = S["Note"];
export type SearchResponse = S["SearchResponse"];
export type PerfumeResponse = S["PerfumeResponse"];
export type NoteMatchResponse = S["NoteMatchResponse"];
export type NoteResponse = S["NoteResponse"];
export type FamiliesResponse = S["FamiliesResponse"];
export type TodayResponse = S["TodayResponse"];
export type ApiError = S["ErrorResponse"];

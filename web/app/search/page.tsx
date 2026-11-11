import type { Metadata } from "next";
import { PerfumeCard } from "@/components/PerfumeCard";
import { SearchBar } from "@/components/SearchBar";
import { ApiRequestError, serverGet } from "@/lib/api";
import type { SearchResponse } from "@/lib/types";

type Props = { searchParams: Promise<{ q?: string }> };

export async function generateMetadata({ searchParams }: Props): Promise<Metadata> {
  const { q } = await searchParams;
  return { title: q ? `“${q}”` : "Search" };
}

export default async function SearchPage({ searchParams }: Props) {
  const q = ((await searchParams).q ?? "").trim();
  let data: SearchResponse | null = null;
  let error = "";

  if (q.length >= 2) {
    try {
      data = await serverGet<SearchResponse>(`/api/v1/perfumes/search?q=${encodeURIComponent(q)}`);
    } catch (e) {
      error = e instanceof ApiRequestError ? e.message : "The perfume database is unavailable right now.";
    }
  }

  return (
    <div className="pt-10">
      <div className="max-w-2xl">
        <SearchBar initial={q} />
      </div>
      {error && <p className="mt-8 text-muted">{error}</p>}
      {data && (
        <>
          <p className="mt-8 text-sm text-faint">
            {data.total === 0
              ? `No perfumes match “${data.query}”.`
              : `${data.total} perfume${data.total === 1 ? "" : "s"} match “${data.query}”${data.total > data.results.length ? ` · showing the first ${data.results.length}` : ""}`}
          </p>
          <ul className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
            {data.results.map((p) => (
              <li key={p.id}>
                <PerfumeCard perfume={p} />
              </li>
            ))}
          </ul>
        </>
      )}
    </div>
  );
}

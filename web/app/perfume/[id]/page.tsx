import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import { Bottle, BottleStage } from "@/components/Bottle";
import { FamilyBadge } from "@/components/FamilyBadge";
import { PerfumeCard } from "@/components/PerfumeCard";
import { NoteList, Pyramid } from "@/components/Pyramid";
import { ShelfButtons } from "@/components/ShelfButtons";
import { ApiRequestError, serverGet } from "@/lib/api";
import type { PerfumeResponse } from "@/lib/types";

type Props = { params: Promise<{ id: string }> };

async function load(id: string): Promise<PerfumeResponse> {
  try {
    return await serverGet<PerfumeResponse>(`/api/v1/perfumes/${encodeURIComponent(id)}`);
  } catch (e) {
    if (e instanceof ApiRequestError && (e.status === 404 || e.status === 400)) notFound();
    throw e;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  try {
    const { perfume } = await load((await params).id);
    return { title: `${perfume.name} by ${perfume.house}`, description: perfume.character };
  } catch {
    return { title: "Perfume" };
  }
}

const SEASON_ICON: Record<string, string> = { spring: "🌱", summer: "☀️", autumn: "🍂", winter: "❄️" };
const HOUR_ICON: Record<string, string> = { day: "🌤️", evening: "🌆", night: "🌙" };

export default async function PerfumePage({ params }: Props) {
  const { perfume: p, neighbours } = await load((await params).id);
  const facts = [
    p.concentration && { k: "Concentration", v: p.concentration },
    p.year && { k: "Launched", v: String(p.year) },
    p.perfumer && { k: "Perfumer", v: p.perfumer },
    p.subfamily && { k: "Subfamily", v: p.subfamily },
  ].filter(Boolean) as { k: string; v: string }[];

  return (
    <div className="pt-8">
      <Link href="/" className="text-sm text-faint hover:text-accent">
        ← Notewise
      </Link>

      <section className="mt-6 grid gap-10 md:grid-cols-[minmax(0,1fr)_minmax(0,1.3fr)]">
        <div className="flex flex-col items-center gap-6 self-start overflow-hidden rounded-3xl border border-line bg-card pb-8 md:sticky md:top-6">
          <BottleStage familyId={p.family.id} className="w-full pb-6 pt-10">
            <Bottle id={p.id} name={p.name} house={p.house} familyId={p.family.id} concentration={p.concentration} size={210} animate />
          </BottleStage>
          <div className="flex flex-wrap justify-center gap-2">
            {p.seasons.map((s) => (
              <span key={s} className="rounded-full bg-paper px-3 py-1 text-sm capitalize text-muted">
                <span aria-hidden>{SEASON_ICON[s] ?? "•"}</span> {s}
              </span>
            ))}
            {p.hours.map((h) => (
              <span key={h} className="rounded-full bg-paper px-3 py-1 text-sm capitalize text-muted">
                <span aria-hidden>{HOUR_ICON[h] ?? "•"}</span> {h}
              </span>
            ))}
          </div>
        </div>

        <div>
          <p className="text-sm uppercase tracking-[0.18em] text-faint">{p.house}</p>
          <h1 className="mt-1 font-display text-5xl leading-tight tracking-tight">{p.name}</h1>
          <div className="mt-3">
            <FamilyBadge family={p.family} size="md" />
          </div>
          {p.character && <p className="mt-5 text-lg text-muted">{p.character}</p>}
          {p.characterWords.length > 0 && (
            <p className="mt-3 text-sm text-faint">
              Feels {p.characterWords.join(", ")}.
            </p>
          )}

          <div className="mt-6">
            <ShelfButtons perfume={p} />
          </div>

          {facts.length > 0 && (
            <dl className="mt-8 grid grid-cols-2 gap-x-6 gap-y-3 border-t border-line pt-6 text-sm">
              {facts.map((f) => (
                <div key={f.k}>
                  <dt className="text-faint">{f.k}</dt>
                  <dd className="first-letter:uppercase">{f.v}</dd>
                </div>
              ))}
            </dl>
          )}

          <div className="mt-8 border-t border-line pt-6">
            <h2 className="font-display text-2xl">{p.pyramid ? "Note pyramid" : "Notes"}</h2>
            <div className="mt-4">
              {p.pyramid ? <Pyramid pyramid={p.pyramid} familyId={p.family.id} /> : <NoteList notes={p.notes} />}
            </div>
            {p.signatureNotes.length > 0 && (
              <p className="mt-4 text-sm text-muted">
                Signature notes:{" "}
                {p.signatureNotes.map((n, i) => (
                  <span key={n}>
                    {i > 0 && ", "}
                    <Link className="underline decoration-line underline-offset-2 hover:text-accent" href={`/notes?notes=${encodeURIComponent(n)}`}>
                      {n}
                    </Link>
                  </span>
                ))}
              </p>
            )}
          </div>
        </div>
      </section>

      {neighbours.length > 0 && (
        <section className="mt-16" aria-labelledby="similar">
          <h2 id="similar" className="font-display text-3xl">If you like {p.name}</h2>
          <p className="mt-1 text-muted">Perfumes that share its structure, and how they differ.</p>
          <ul className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-4">
            {neighbours.map((n) => (
              <li key={n.id}>
                <PerfumeCard
                  perfume={n}
                  footer={
                    <>
                      <p>
                        <span className="text-faint">Shares </span>
                        {n.sharedNotes.join(", ").toLowerCase()}
                      </p>
                      {n.difference && <p className="mt-1 text-faint">{n.difference}</p>}
                    </>
                  }
                />
              </li>
            ))}
          </ul>
        </section>
      )}

      {p.sourceUrl && (
        <p className="mt-12 text-xs text-faint">
          Notes as published by the house or an authorized retailer, via{" "}
          <a href={p.sourceUrl} target="_blank" rel="noreferrer" className="underline underline-offset-2 hover:text-muted">
            Scento
          </a>
          .
        </p>
      )}
    </div>
  );
}

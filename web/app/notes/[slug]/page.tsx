import type { Metadata } from "next";
import Image from "next/image";
import Link from "next/link";
import { notFound } from "next/navigation";
import { FamilyBadge } from "@/components/FamilyBadge";
import { NoteChip } from "@/components/NoteChip";
import { PerfumeCard } from "@/components/PerfumeCard";
import { ApiRequestError, serverGet } from "@/lib/api";
import { notePhoto } from "@/lib/notes";
import type { NoteResponse } from "@/lib/types";

type Props = { params: Promise<{ slug: string }> };

async function load(slug: string): Promise<NoteResponse> {
  try {
    return await serverGet<NoteResponse>(`/api/v1/notes/${encodeURIComponent(slug)}`);
  } catch (e) {
    if (e instanceof ApiRequestError && (e.status === 404 || e.status === 400)) notFound();
    throw e;
  }
}

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  try {
    const { note } = await load((await params).slug);
    return { title: `${note.name} in perfume`, description: note.smellsLike };
  } catch {
    return { title: "Note" };
  }
}

const TIER_TEXT: Record<string, string> = {
  top: "In the opening",
  heart: "In the heart",
  base: "In the dry-down",
  listed: "Listed among its notes",
};

export default async function NotePage({ params }: Props) {
  const { note, perfumes, total } = await load((await params).slug);
  const photo = notePhoto(note.name);
  const facts = [
    note.kind && { k: "Natural or synthetic", v: note.kind },
    note.category && { k: "Category", v: note.category },
    note.origin && { k: "Origin", v: note.origin },
  ].filter(Boolean) as { k: string; v: string }[];

  return (
    <div className="pt-8">
      <Link href="/ingredients" className="text-sm text-faint hover:text-accent">
        ← All ingredients
      </Link>

      <section className="mt-6 grid items-center gap-8 md:grid-cols-[auto_1fr]">
        <figure className="mx-auto">
          <div className="relative h-56 w-56 overflow-hidden rounded-full bg-accent-soft ring-1 ring-line sm:h-64 sm:w-64">
            {photo ? (
              <Image src={photo.src} alt={note.name} fill sizes="256px" className="object-cover" priority />
            ) : (
              <span className="flex h-full w-full items-center justify-center font-display text-7xl text-accent">{note.name[0]}</span>
            )}
          </div>
          {photo && (
            <figcaption className="mt-3 max-w-64 text-center text-[11px] leading-snug text-faint">
              Photo: {photo.credit.author} ·{" "}
              <a href={photo.credit.licenseUrl || photo.credit.source} target="_blank" rel="noreferrer" className="underline underline-offset-2">
                {photo.credit.license}
              </a>
            </figcaption>
          )}
        </figure>

        <div>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Ingredient</p>
          <h1 className="mt-1 font-display text-5xl tracking-tight">{note.name}</h1>
          {note.family?.name && (
            <div className="mt-3">
              <FamilyBadge family={note.family} size="md" />
            </div>
          )}
          {note.smellsLike && <p className="mt-5 max-w-2xl text-xl text-muted">{note.smellsLike}</p>}
          {facts.length > 0 && (
            <dl className="mt-6 flex flex-wrap gap-x-10 gap-y-3 text-sm">
              {facts.map((f) => (
                <div key={f.k}>
                  <dt className="text-faint">{f.k}</dt>
                  <dd className="first-letter:uppercase">{f.v}</dd>
                </div>
              ))}
            </dl>
          )}
          {note.pairsWith.length > 0 && (
            <div className="mt-6">
              <p className="text-sm text-muted">Pairs well with</p>
              <div className="mt-2 flex flex-wrap gap-2">
                {note.pairsWith.map((n) => (
                  <NoteChip key={n} name={n} size="sm" />
                ))}
              </div>
            </div>
          )}
          {note.aliases.length > 0 && (
            <p className="mt-5 text-sm text-faint">Also listed as: {note.aliases.slice(0, 6).join(", ")}.</p>
          )}
        </div>
      </section>

      {perfumes.length > 0 && (
        <section className="mt-16" aria-labelledby="with-note">
          <h2 id="with-note" className="font-display text-3xl">
            Perfumes with {note.name.toLowerCase()}
          </h2>
          <p className="mt-1 text-muted">
            {total} perfume{total === 1 ? "" : "s"} in the database list it
            {total > perfumes.length ? `; here are the first ${perfumes.length}` : ""}.{" "}
            <Link href={`/notes?notes=${encodeURIComponent(note.name)}`} className="underline decoration-line underline-offset-2 hover:text-accent">
              Combine it with other notes
            </Link>
            .
          </p>
          <ul className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
            {perfumes.map((p) => (
              <li key={p.id}>
                <PerfumeCard
                  perfume={p}
                  compact
                  footer={
                    Object.values(p.tiers)[0] && (
                      <span className="text-xs">{TIER_TEXT[Object.values(p.tiers)[0]] ?? Object.values(p.tiers)[0]}</span>
                    )
                  }
                />
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}

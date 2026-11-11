import Link from "next/link";
import { Bottle } from "@/components/Bottle";
import { SearchBar } from "@/components/SearchBar";
import { TodayPick } from "@/components/TodayPick";
import { serverGet } from "@/lib/api";
import { familyColor } from "@/lib/families";
import type { FamiliesResponse } from "@/lib/types";

const FEATURED = [
  { id: "eau-sauvage-dior", name: "Eau Sauvage", house: "Dior", familyId: "sunlit-citrus", concentration: "eau de toilette" },
  { id: "layton-parfums-de-marly", name: "Layton", house: "Parfums de Marly", familyId: "wood-root", concentration: "eau de parfum" },
  { id: "hypnotic-poison-dior", name: "Hypnotic Poison", house: "Dior", familyId: "velvet-sweet", concentration: "eau de toilette" },
];

async function families() {
  try {
    return (await serverGet<FamiliesResponse>("/api/v1/families")).families;
  } catch {
    return []; // the page still works if the API is waking up
  }
}

export default async function Home() {
  const fams = await families();

  return (
    <div className="flex flex-col gap-16 pt-10 sm:pt-16">
      <section className="grid items-center gap-10 md:grid-cols-[1.2fr_1fr]">
        <div>
          <p className="rise text-xs uppercase tracking-[0.2em] text-faint">A field guide to perfume</p>
          <h1 className="rise rise-1 mt-3 font-display text-5xl leading-[1.05] tracking-tight sm:text-6xl">
            Read a perfume <span className="italic text-accent">by its notes.</span>
          </h1>
          <p className="rise rise-2 mt-5 max-w-lg text-lg text-muted">
            Explore note pyramids, find scents built on the notes you love, get a pick for today&apos;s weather and see
            what your shelf says about your taste.
          </p>
          <div className="rise rise-3 mt-8 max-w-xl">
            <SearchBar />
            <p className="mt-3 text-sm text-faint">
              Or browse by note in the{" "}
              <Link href="/notes" className="text-muted underline decoration-line underline-offset-2 hover:text-accent">
                note finder
              </Link>
              .
            </p>
          </div>
        </div>
        <div className="flex items-end justify-center gap-2 sm:gap-4" aria-hidden>
          {FEATURED.map((f, i) => (
            <Link key={f.id} href={`/perfume/${f.id}`} tabIndex={-1} className={`rise rise-${i + 1} transition hover:-translate-y-1`}>
              <Bottle {...f} size={i === 1 ? 150 : 112} animate className={i === 1 ? "h-auto w-[34vw] max-w-[150px]" : "h-auto w-[25vw] max-w-[112px]"} />
            </Link>
          ))}
        </div>
      </section>

      <TodayPick />

      {fams.length > 0 && (
        <section aria-labelledby="families">
          <div className="flex items-end justify-between gap-4">
            <div>
              <p className="text-xs uppercase tracking-[0.18em] text-faint">The map</p>
              <h2 id="families" className="mt-1 font-display text-3xl">{fams.length} olfactory families</h2>
            </div>
          </div>
          <ul className="mt-6 grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {fams.map((f) => (
              <li key={f.id} className="flex gap-3 rounded-2xl border border-line bg-card p-4">
                <span aria-hidden className="mt-1 h-8 w-2 shrink-0 rounded-full" style={{ background: familyColor(f.id) }} />
                <div>
                  <p className="font-display text-lg">
                    {f.name} <span className="font-sans text-xs text-faint">· {f.count} perfumes</span>
                  </p>
                  <p className="mt-0.5 text-sm text-muted">{f.definition}</p>
                </div>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}

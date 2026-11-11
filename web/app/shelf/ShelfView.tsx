"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { BarChart, Coverage } from "@/components/BarChart";
import { PerfumeCard } from "@/components/PerfumeCard";
import { clientFetch } from "@/lib/api";
import { FAMILY_COLORS } from "@/lib/families";
import { useShelf, type ShelfList } from "@/lib/shelf";
import type { ProfileResponse } from "@/lib/types";

const SEASONS = { spring: "🌱", summer: "☀️", autumn: "🍂", winter: "❄️" };
const HOURS = { day: "🌤️", evening: "🌆", night: "🌙" };
// "Wood & Root" → "wood-root", the id the colors are keyed by.
function swatchFor(label: string) {
  return FAMILY_COLORS[label.toLowerCase().replace(/&/g, "").replace(/\W+/g, "-")];
}

export function ShelfView() {
  const { shelf, ready, toggle } = useShelf();
  const [tab, setTab] = useState<ShelfList>("owned");
  const [profile, setProfile] = useState<ProfileResponse | null>(null);
  const [state, setState] = useState<"idle" | "loading" | "error">("idle");

  const ownedKey = shelf.owned.map((p) => p.id).join(",");
  useEffect(() => {
    if (!ready || !shelf.owned.length) {
      setProfile(null);
      return;
    }
    let cancelled = false;
    setState("loading");
    clientFetch<ProfileResponse>("/api/v1/collection/profile", {
      method: "POST",
      body: JSON.stringify({ perfumeIds: shelf.owned.map((p) => p.id) }),
    })
      .then((p) => !cancelled && (setProfile(p), setState("idle")))
      .catch(() => !cancelled && setState("error"));
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [ready, ownedKey]);

  if (!ready) return <div className="mt-10 h-40 animate-pulse rounded-3xl bg-card" />;

  const list = shelf[tab];

  return (
    <div className="mt-10 flex flex-col gap-12">
      <section>
        <div className="inline-flex rounded-full border border-line p-1" role="tablist">
          {(["owned", "wishlist"] as const).map((t) => (
            <button
              key={t}
              role="tab"
              aria-selected={tab === t}
              onClick={() => setTab(t)}
              className={`rounded-full px-4 py-1.5 text-sm transition ${tab === t ? "bg-ink text-paper" : "text-muted hover:text-ink"}`}
            >
              {t === "owned" ? "I own" : "Wishlist"} · {shelf[t].length}
            </button>
          ))}
        </div>

        {list.length === 0 ? (
          <div className="mt-6 rounded-3xl border border-dashed border-line p-8 text-center">
            <p className="font-display text-2xl">{tab === "owned" ? "Your shelf is empty" : "Nothing on your wishlist yet"}</p>
            <p className="mt-2 text-muted">
              <Link href="/" className="underline decoration-line underline-offset-2 hover:text-accent">
                Search
              </Link>{" "}
              for a perfume and tap “{tab === "owned" ? "I own this" : "Wishlist"}”.
            </p>
          </div>
        ) : (
          <ul className="mt-6 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
            {list.map((p) => (
              <li key={p.id}>
                <PerfumeCard
                  perfume={p}
                  compact
                  footer={
                    <button type="button" onClick={() => toggle(tab, p)} className="text-xs text-faint hover:text-accent">
                      Remove
                    </button>
                  }
                />
              </li>
            ))}
          </ul>
        )}
      </section>

      {shelf.owned.length > 0 && (
        <section aria-labelledby="profile">
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Scent profile</p>
          {state === "loading" && !profile && <p className="mt-2 text-faint">Reading your shelf…</p>}
          {state === "error" && <p className="mt-2 text-muted">We couldn&apos;t build your profile right now. Please try again later.</p>}
          {profile && (
            <div className="rise">
              <h2 id="profile" className="mt-1 max-w-3xl font-display text-4xl leading-tight">
                {profile.headline}
              </h2>
              {profile.characterWords.length > 0 && (
                <p className="mt-3 text-muted">
                  In one breath: {profile.characterWords.slice(0, 5).map((w) => w.label).join(", ")}.
                </p>
              )}
              <div className="mt-6 grid gap-4 lg:grid-cols-2">
                <BarChart title="Families" data={profile.families} swatch={swatchFor} />
                <BarChart title="Most frequent notes" data={profile.topNotes} />
                <Coverage title="Seasons covered" data={profile.seasons} icons={SEASONS} />
                <Coverage title="Moments covered" data={profile.hours} icons={HOURS} />
              </div>
              {profile.gaps.length > 0 && (
                <p className="mt-4 text-muted">
                  Missing from your wardrobe: <strong className="font-medium text-ink">{profile.gaps.join(", ")}</strong>.
                </p>
              )}
              {profile.missing.length > 0 && (
                <p className="mt-2 text-sm text-faint">{profile.missing.length} perfume(s) could not be found in the database.</p>
              )}

              {profile.suggestions.length > 0 && (
                <div className="mt-10">
                  <h3 className="font-display text-3xl">Round out your wardrobe</h3>
                  <p className="mt-1 text-muted">Built from what you already own, one scent per moment.</p>
                  <ul className="mt-5 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
                    {profile.suggestions.map((s) => (
                      <li key={s.id}>
                        <PerfumeCard
                          perfume={s}
                          compact
                          footer={
                            <>
                              {s.role && <p className="text-xs uppercase tracking-wider text-accent">{s.role}</p>}
                              <p className="mt-1 line-clamp-3">{s.why}</p>
                            </>
                          }
                        />
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </section>
      )}
    </div>
  );
}

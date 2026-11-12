"use client";

import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { PerfumeCard } from "@/components/PerfumeCard";
import { ApiRequestError, clientFetch } from "@/lib/api";
import { POPULAR_NOTES } from "@/lib/families";
import Image from "next/image";
import { notePhoto } from "@/lib/notes";
import type { NoteMatchResponse } from "@/lib/types";

const TIERS = [
  { id: "any", label: "Anywhere" },
  { id: "top", label: "Opening" },
  { id: "heart", label: "Heart" },
  { id: "base", label: "Dry-down" },
];

export function NoteFinder() {
  const router = useRouter();
  const params = useSearchParams();
  const notes = (params.get("notes") ?? "").split(",").map((n) => n.trim()).filter(Boolean);
  const tier = params.get("tier") ?? "any";
  const [input, setInput] = useState("");
  const [data, setData] = useState<NoteMatchResponse | null>(null);
  const [state, setState] = useState<"idle" | "loading" | "error">("idle");
  const [error, setError] = useState("");

  // The URL is the state: results are shareable and survive a refresh.
  function update(nextNotes: string[], nextTier = tier) {
    const q = new URLSearchParams();
    if (nextNotes.length) q.set("notes", nextNotes.join(","));
    if (nextTier !== "any") q.set("tier", nextTier);
    router.replace(`/notes${q.size ? `?${q}` : ""}`, { scroll: false });
  }

  function add(note: string) {
    const n = note.trim();
    if (!n || notes.length >= 4 || notes.some((x) => x.toLowerCase() === n.toLowerCase())) return;
    update([...notes, n]);
    setInput("");
  }

  const key = `${notes.join(",")}|${tier}`;
  useEffect(() => {
    if (!notes.length) {
      setData(null);
      return;
    }
    let cancelled = false;
    setState("loading");
    const q = new URLSearchParams({ notes: notes.join(","), tier });
    clientFetch<NoteMatchResponse>(`/api/v1/notes/match?${q}`)
      .then((d) => !cancelled && (setData(d), setState("idle")))
      .catch((e) => {
        if (cancelled) return;
        setError(e instanceof ApiRequestError && e.status < 500 ? e.message : "The perfume database is unavailable right now.");
        setState("error");
      });
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [key]);

  const suggestions = POPULAR_NOTES.filter((n) => !notes.some((x) => x.toLowerCase() === n.toLowerCase()));

  return (
    <div className="mt-8">
      <div className="rounded-3xl border border-line bg-card p-5 sm:p-6">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            add(input);
          }}
          className="flex flex-wrap items-center gap-2"
        >
          {notes.map((n) => (
            <span key={n} className="inline-flex items-center gap-1.5 rounded-full bg-accent-soft py-1 pl-1 pr-1 text-sm text-accent">
              <Thumb name={n} />
              {n}
              <button
                type="button"
                onClick={() => update(notes.filter((x) => x !== n))}
                aria-label={`Remove ${n}`}
                className="rounded-full px-2 hover:bg-card"
              >
                ×
              </button>
            </span>
          ))}
          <label htmlFor="note" className="sr-only">
            Add a note
          </label>
          <input
            id="note"
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={notes.length >= 4}
            placeholder={notes.length >= 4 ? "Four notes max" : "Type a note and press Enter"}
            className="min-w-48 flex-1 bg-transparent py-1.5 outline-none placeholder:text-faint"
            maxLength={40}
          />
        </form>

        <div className="mt-4 flex flex-wrap gap-1.5">
          {suggestions.slice(0, 18).map((n) => (
            <button
              key={n}
              type="button"
              onClick={() => add(n)}
              disabled={notes.length >= 4}
              className="inline-flex items-center gap-1.5 rounded-full border border-line py-1 pl-1 pr-3 text-sm text-muted transition hover:border-accent hover:text-ink disabled:opacity-40"
            >
              <Thumb name={n} /> {n}
            </button>
          ))}
        </div>

        <div className="mt-5 inline-flex rounded-full border border-line p-1" role="radiogroup" aria-label="Where the notes sit">
          {TIERS.map((t) => (
            <button
              key={t.id}
              type="button"
              role="radio"
              aria-checked={tier === t.id}
              onClick={() => update(notes, t.id)}
              className={`rounded-full px-4 py-1.5 text-sm transition ${tier === t.id ? "bg-ink text-paper" : "text-muted hover:text-ink"}`}
            >
              {t.label}
            </button>
          ))}
        </div>
      </div>

      <div aria-live="polite" className="mt-8">
        {!notes.length && <p className="text-muted">Add a note to start, for example Vanilla or Bergamot.</p>}
        {state === "error" && <p className="text-muted">{error}</p>}
        {state === "loading" && <p className="text-faint">Searching…</p>}
        {state === "idle" && data && (
          <>
            <p className="text-sm text-faint">
              {data.total === 0
                ? "No perfumes contain all of these notes. Try removing one."
                : `${data.total} perfume${data.total === 1 ? "" : "s"} contain ${data.notes.join(" + ")}${data.total > data.results.length ? ` · top ${data.results.length}` : ""}`}
            </p>
            <ul className="mt-4 grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-5">
              {data.results.map((p) => (
                <li key={p.id}>
                  <PerfumeCard
                    perfume={p}
                    compact
                    footer={
                      <ul className="flex flex-wrap gap-1">
                        {Object.entries(p.tiers).map(([note, where]) => (
                          <li key={note} className="rounded bg-paper px-1.5 py-0.5 text-xs">
                            {note} <span className="text-faint">· {where}</span>
                          </li>
                        ))}
                      </ul>
                    }
                  />
                </li>
              ))}
            </ul>
          </>
        )}
      </div>
    </div>
  );
}

function Thumb({ name }: { name: string }) {
  const photo = notePhoto(name);
  return photo ? (
    <Image src={photo.src} alt="" width={48} height={48} className="h-6 w-6 rounded-full object-cover" />
  ) : (
    <span aria-hidden className="flex h-6 w-6 items-center justify-center rounded-full bg-card text-xs">
      {name[0]}
    </span>
  );
}

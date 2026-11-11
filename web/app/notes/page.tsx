import type { Metadata } from "next";
import { Suspense } from "react";
import { NoteFinder } from "./NoteFinder";

export const metadata: Metadata = { title: "Note finder" };

export default function NotesPage() {
  return (
    <div className="pt-10">
      <p className="text-xs uppercase tracking-[0.18em] text-faint">Note finder</p>
      <h1 className="mt-1 font-display text-5xl tracking-tight">Start from the notes you love</h1>
      <p className="mt-3 max-w-2xl text-lg text-muted">
        Pick up to four notes and find perfumes that contain all of them, optionally only in the opening, the heart or the
        dry-down.
      </p>
      <Suspense>
        <NoteFinder />
      </Suspense>
    </div>
  );
}

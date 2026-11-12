import type { Metadata } from "next";
import { NoteChip } from "@/components/NoteChip";
import { allNotePhotos } from "@/lib/notes";

export const metadata: Metadata = { title: "Ingredients" };

function label(key: string) {
  return key.replace(/-/g, " ").replace(/\b\w/g, (c) => c.toUpperCase()).replace(/ Of The /, " of the ");
}

export default function IngredientsPage() {
  const notes = allNotePhotos().sort((a, b) => a.key.localeCompare(b.key));
  return (
    <div className="pt-10">
      <p className="text-xs uppercase tracking-[0.18em] text-faint">Ingredients</p>
      <h1 className="mt-1 font-display text-5xl tracking-tight">What perfumes are made of</h1>
      <p className="mt-3 max-w-2xl text-lg text-muted">
        {notes.length} of the notes perfumers use most. Open one to see what it smells like, what it pairs with and which
        perfumes are built on it.
      </p>
      <ul className="mt-10 grid grid-cols-3 gap-x-4 gap-y-8 sm:grid-cols-5 lg:grid-cols-8">
        {notes.map((n) => (
          <li key={n.key} className="flex justify-center">
            <NoteChip name={label(n.key)} size="lg" />
          </li>
        ))}
      </ul>
    </div>
  );
}

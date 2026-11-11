"use client";

import { useRouter } from "next/navigation";
import { useState } from "react";

export function SearchBar({ initial = "", autoFocus = false }: { initial?: string; autoFocus?: boolean }) {
  const router = useRouter();
  const [q, setQ] = useState(initial);
  const valid = q.trim().length >= 2;

  return (
    <form
      role="search"
      onSubmit={(e) => {
        e.preventDefault();
        if (valid) router.push(`/search?q=${encodeURIComponent(q.trim())}`);
      }}
      className="flex w-full items-center gap-2 rounded-full border border-line bg-card p-1.5 pl-5 shadow-sm focus-within:border-accent"
    >
      <label htmlFor="search" className="sr-only">
        Search perfumes, houses or notes
      </label>
      <input
        id="search"
        value={q}
        onChange={(e) => setQ(e.target.value)}
        autoFocus={autoFocus}
        placeholder="Try “Sauvage”, “Le Labo” or “vanilla”"
        className="min-w-0 flex-1 bg-transparent py-2 text-base outline-none placeholder:text-faint"
        maxLength={80}
      />
      <button
        type="submit"
        disabled={!valid}
        className="rounded-full bg-ink px-5 py-2.5 text-sm font-medium text-paper transition hover:opacity-90 disabled:opacity-40"
      >
        Search
      </button>
    </form>
  );
}

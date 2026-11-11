"use client";

import { useShelf } from "@/lib/shelf";
import type { Summary } from "@/lib/types";

export function ShelfButtons({ perfume }: { perfume: Summary }) {
  const { has, toggle, ready } = useShelf();
  const owned = ready && has("owned", perfume.id);
  const wished = ready && has("wishlist", perfume.id);
  const item: Summary = { id: perfume.id, name: perfume.name, house: perfume.house, family: { id: perfume.family.id, name: perfume.family.name } };

  return (
    <div className="flex flex-wrap gap-2">
      <button
        type="button"
        onClick={() => toggle("owned", item)}
        aria-pressed={owned}
        className={`rounded-full px-5 py-2.5 text-sm font-medium transition ${
          owned ? "bg-accent text-paper" : "bg-ink text-paper hover:opacity-90"
        }`}
      >
        {owned ? "✓ On my shelf" : "I own this"}
      </button>
      <button
        type="button"
        onClick={() => toggle("wishlist", item)}
        aria-pressed={wished}
        disabled={owned}
        className={`rounded-full border px-5 py-2.5 text-sm font-medium transition disabled:opacity-40 ${
          wished ? "border-accent bg-accent-soft text-accent" : "border-line bg-card hover:border-accent"
        }`}
      >
        {wished ? "♥ On my wishlist" : "♡ Wishlist"}
      </button>
    </div>
  );
}

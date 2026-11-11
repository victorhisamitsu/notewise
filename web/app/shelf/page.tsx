import type { Metadata } from "next";
import { ShelfView } from "./ShelfView";

export const metadata: Metadata = { title: "My shelf" };

export default function ShelfPage() {
  return (
    <div className="pt-10">
      <p className="text-xs uppercase tracking-[0.18em] text-faint">My shelf</p>
      <h1 className="mt-1 font-display text-5xl tracking-tight">What your shelf says about you</h1>
      <p className="mt-3 max-w-2xl text-lg text-muted">
        Your shelf is saved in this browser only, with no account needed. Add perfumes you own to see your scent profile and
        what is missing from your wardrobe.
      </p>
      <ShelfView />
    </div>
  );
}

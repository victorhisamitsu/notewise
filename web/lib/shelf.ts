"use client";

import { useCallback, useEffect, useState } from "react";
import type { Summary } from "./types";

// The shelf lives in the visitor's browser: no account needed. Storage can
// be unavailable (private mode, blocked cookies), so every access is guarded
// and the UI works with an empty shelf.

export type ShelfList = "owned" | "wishlist";
type Shelf = Record<ShelfList, Summary[]>;

const KEY = "notewise.shelf.v1";
const EVENT = "notewise:shelf";
const EMPTY: Shelf = { owned: [], wishlist: [] };

function read(): Shelf {
  try {
    const raw = window.localStorage.getItem(KEY);
    if (!raw) return EMPTY;
    const parsed = JSON.parse(raw) as Partial<Shelf>;
    return { owned: parsed.owned ?? [], wishlist: parsed.wishlist ?? [] };
  } catch {
    return EMPTY;
  }
}

function write(shelf: Shelf) {
  try {
    window.localStorage.setItem(KEY, JSON.stringify(shelf));
  } catch {
    /* storage unavailable: keep working in memory for this page */
  }
  window.dispatchEvent(new Event(EVENT));
}

export function useShelf() {
  const [shelf, setShelf] = useState<Shelf>(EMPTY);
  const [ready, setReady] = useState(false);

  useEffect(() => {
    const sync = () => setShelf(read());
    sync();
    setReady(true);
    window.addEventListener(EVENT, sync);
    window.addEventListener("storage", sync); // other tabs
    return () => {
      window.removeEventListener(EVENT, sync);
      window.removeEventListener("storage", sync);
    };
  }, []);

  const has = useCallback((list: ShelfList, id: string) => shelf[list].some((p) => p.id === id), [shelf]);

  const toggle = useCallback((list: ShelfList, item: Summary) => {
    const current = read();
    const exists = current[list].some((p) => p.id === item.id);
    const next: Shelf = {
      ...current,
      [list]: exists ? current[list].filter((p) => p.id !== item.id) : [item, ...current[list]],
    };
    // A perfume you own is no longer on your wishlist.
    if (list === "owned" && !exists) next.wishlist = next.wishlist.filter((p) => p.id !== item.id);
    write(next);
  }, []);

  return { shelf, ready, has, toggle };
}

import type { Pyramid as PyramidT } from "@/lib/types";
import { familyColor, shade } from "@/lib/families";

const TIERS: { key: keyof PyramidT; label: string; hint: string }[] = [
  { key: "top", label: "Top", hint: "first minutes" },
  { key: "heart", label: "Heart", hint: "the core, after ~30 min" },
  { key: "base", label: "Base", hint: "the dry-down, for hours" },
];

/** The olfactory pyramid: narrow and light at the top, wide and deep at the base. */
export function Pyramid({ pyramid, familyId }: { pyramid: PyramidT; familyId: string }) {
  const color = familyColor(familyId);
  const tints = [shade(color, 0.72), shade(color, 0.55), shade(color, 0.38)];
  return (
    <ol className="flex flex-col items-center gap-2" aria-label="Note pyramid">
      {TIERS.map((t, i) => (
        <li
          key={t.key}
          className="rise w-full rounded-2xl px-4 py-3 text-center text-[#2a2420]"
          style={{ maxWidth: `${62 + i * 19}%`, background: tints[i], animationDelay: `${0.1 + i * 0.12}s` }}
        >
          <p className="text-[11px] uppercase tracking-[0.18em] opacity-70">
            {t.label} <span className="normal-case tracking-normal">· {t.hint}</span>
          </p>
          <p className="mt-1 font-display text-base leading-snug">
            {pyramid[t.key].length ? pyramid[t.key].join(" · ") : "—"}
          </p>
        </li>
      ))}
    </ol>
  );
}

export function NoteList({ notes }: { notes: string[] }) {
  return (
    <ul className="flex flex-wrap gap-2" aria-label="Notes">
      {notes.map((n) => (
        <li key={n} className="rounded-full border border-line bg-card px-3 py-1 text-sm">
          {n}
        </li>
      ))}
    </ul>
  );
}

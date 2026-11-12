import type { Pyramid as PyramidT } from "@/lib/types";
import { familyColor, shade } from "@/lib/families";
import { NoteChip } from "./NoteChip";

const TIERS: { key: keyof PyramidT; label: string; hint: string }[] = [
  { key: "top", label: "Top", hint: "first minutes" },
  { key: "heart", label: "Heart", hint: "the core, after ~30 min" },
  { key: "base", label: "Base", hint: "the dry-down, for hours" },
];

/** The olfactory pyramid: each tier lists its notes as ingredient photos. */
export function Pyramid({ pyramid, familyId }: { pyramid: PyramidT; familyId: string }) {
  const color = familyColor(familyId);
  const tints = [shade(color, 0.82), shade(color, 0.72), shade(color, 0.62)];
  return (
    <ol className="flex flex-col gap-2" aria-label="Note pyramid">
      {TIERS.map((t, i) => (
        <li
          key={t.key}
          className="rise rounded-2xl px-4 py-4 text-[#2a2420]"
          style={{ background: tints[i], animationDelay: `${0.1 + i * 0.12}s` }}
        >
          <p className="text-center text-[11px] uppercase tracking-[0.18em] opacity-70">
            {t.label} <span className="normal-case tracking-normal">· {t.hint}</span>
          </p>
          <div className="mt-3 flex flex-wrap justify-center gap-x-2 gap-y-3">
            {pyramid[t.key].length ? pyramid[t.key].map((n) => <NoteChip key={n} name={n} />) : <span>—</span>}
          </div>
        </li>
      ))}
    </ol>
  );
}

export function NoteList({ notes }: { notes: string[] }) {
  return (
    <ul className="flex flex-wrap gap-x-2 gap-y-3" aria-label="Notes">
      {notes.map((n) => (
        <li key={n}>
          <NoteChip name={n} />
        </li>
      ))}
    </ul>
  );
}

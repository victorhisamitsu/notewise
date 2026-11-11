"use client";

import { useState } from "react";
import type { Count } from "@/lib/types";

type Props = {
  title: string;
  data: Count[];
  unit?: string;
  swatch?: (label: string) => string | undefined; // identity dot beside the label
};

/**
 * Single-series horizontal bar chart. One hue (validated --bar token),
 * 14px bars with a 4px rounded data-end, value at the tip in text ink,
 * per-bar hover tooltip and a data table for screen readers.
 */
export function BarChart({ title, data, unit = "perfume", swatch }: Props) {
  const [hover, setHover] = useState<number | null>(null);
  const max = Math.max(1, ...data.map((d) => d.count));
  const id = title.toLowerCase().replace(/\W+/g, "-");

  return (
    <figure className="rounded-2xl border border-line bg-card p-5">
      <figcaption id={id} className="font-display text-lg">
        {title}
      </figcaption>
      <ul className="mt-4 flex flex-col gap-2.5" aria-hidden>
        {data.map((d, i) => (
          <li
            key={d.label}
            className="relative grid grid-cols-[7.5rem_1fr] items-center gap-3 text-sm"
            onMouseEnter={() => setHover(i)}
            onMouseLeave={() => setHover(null)}
          >
            <span className="flex min-w-0 items-center gap-1.5 text-muted">
              {swatch?.(d.label) && <span className="h-2 w-2 shrink-0 rounded-full" style={{ background: swatch(d.label) }} />}
              <span className="truncate capitalize" title={d.label}>
                {d.label}
              </span>
            </span>
            <span className="flex items-center gap-2 py-1">
              <span
                className="h-3.5 rounded-r bg-bar transition-opacity"
                style={{ width: `${(d.count / max) * 85}%`, minWidth: d.count ? 4 : 0, opacity: hover === null || hover === i ? 1 : 0.45 }}
              />
              <span className="tabular-nums text-ink">{d.count}</span>
            </span>
            {hover === i && (
              <span className="pointer-events-none absolute -top-8 left-32 z-10 whitespace-nowrap rounded-md bg-ink px-2 py-1 text-xs text-paper shadow">
                {d.label}: {d.count} {unit}
                {d.count === 1 ? "" : "s"}
              </span>
            )}
          </li>
        ))}
      </ul>
      <table className="sr-only" aria-labelledby={id}>
        <thead>
          <tr>
            <th>Label</th>
            <th>Count</th>
          </tr>
        </thead>
        <tbody>
          {data.map((d) => (
            <tr key={d.label}>
              <td>{d.label}</td>
              <td>{d.count}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </figure>
  );
}

/** Coverage cells for a fixed set (seasons, moments): filled when covered. */
export function Coverage({ title, data, icons }: { title: string; data: Count[]; icons: Record<string, string> }) {
  return (
    <figure className="rounded-2xl border border-line bg-card p-5">
      <figcaption className="font-display text-lg">{title}</figcaption>
      <ul className="mt-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
        {data.map((d) => (
          <li
            key={d.label}
            className={`rounded-xl border px-3 py-3 text-center text-sm ${d.count ? "border-transparent bg-accent-soft" : "border-dashed border-line text-faint"}`}
          >
            <span aria-hidden className="block text-xl">
              {icons[d.label] ?? "•"}
            </span>
            <span className="capitalize">{d.label}</span>
            <span className="block text-xs text-muted">{d.count ? `${d.count} perfume${d.count === 1 ? "" : "s"}` : "not covered"}</span>
          </li>
        ))}
      </ul>
    </figure>
  );
}

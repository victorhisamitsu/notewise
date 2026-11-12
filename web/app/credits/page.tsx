import type { Metadata } from "next";
import Image from "next/image";
import { allNotePhotos } from "@/lib/notes";

export const metadata: Metadata = { title: "Credits" };

export default function CreditsPage() {
  const photos = allNotePhotos().sort((a, b) => a.key.localeCompare(b.key));
  return (
    <div className="pt-10">
      <h1 className="font-display text-5xl tracking-tight">Credits</h1>
      <p className="mt-3 max-w-2xl text-muted">
        Perfume data comes from the{" "}
        <a className="underline underline-offset-2" href="https://www.scento.com/developers/perfume-api" target="_blank" rel="noreferrer">
          Scento Perfume API
        </a>
        , weather from{" "}
        <a className="underline underline-offset-2" href="https://open-meteo.com" target="_blank" rel="noreferrer">
          Open-Meteo
        </a>
        . Ingredient photos are free-licensed works from Wikimedia Commons, cropped to a square. Thank you to their authors.
      </p>
      <ul className="mt-8 grid gap-3 sm:grid-cols-2">
        {photos.map((p) => (
          <li key={p.key} className="flex items-center gap-3 rounded-2xl border border-line bg-card p-3 text-sm">
            <Image src={`/notes/${p.key}.webp`} alt="" width={48} height={48} className="h-12 w-12 rounded-full object-cover" />
            <div className="min-w-0">
              <p className="font-medium capitalize">{p.key.replace(/-/g, " ")}</p>
              <p className="truncate text-faint">
                <a href={p.source} target="_blank" rel="noreferrer" className="underline underline-offset-2 hover:text-muted">
                  {p.title}
                </a>{" "}
                · {p.author} ·{" "}
                {p.licenseUrl ? (
                  <a href={p.licenseUrl} target="_blank" rel="noreferrer" className="underline underline-offset-2 hover:text-muted">
                    {p.license}
                  </a>
                ) : (
                  p.license
                )}
              </p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

import Image from "next/image";
import Link from "next/link";
import { notePhoto, slugify } from "@/lib/notes";

type Props = { name: string; size?: "sm" | "md" | "lg"; caption?: string; link?: boolean };

const PX = { sm: 28, md: 52, lg: 72 } as const;

/** A note as a round ingredient photo + its name; falls back to an initial when there is no photo. */
export function NoteChip({ name, size = "md", caption, link = true }: Props) {
  const photo = notePhoto(name);
  const px = PX[size];
  const inner = (
    <>
      <span
        className="relative inline-flex shrink-0 items-center justify-center overflow-hidden rounded-full bg-accent-soft font-display text-accent ring-1 ring-line"
        style={{ width: px, height: px, fontSize: px * 0.4 }}
      >
        {photo ? (
          <Image src={photo.src} alt="" width={px * 2} height={px * 2} className="h-full w-full object-cover" />
        ) : (
          <span aria-hidden>{name.trim()[0]?.toUpperCase()}</span>
        )}
      </span>
      <span className={size === "sm" ? "text-sm" : "text-center text-sm leading-tight"}>
        {name}
        {caption && <span className="block text-xs text-faint">{caption}</span>}
      </span>
    </>
  );
  const cls =
    size === "sm"
      ? "inline-flex items-center gap-2 rounded-full border border-line bg-card py-1 pl-1 pr-3"
      : "inline-flex w-20 flex-col items-center gap-1.5";
  return link ? (
    <Link href={`/notes/${slugify(name)}`} className={`${cls} transition hover:text-accent`}>
      {inner}
    </Link>
  ) : (
    <span className={cls}>{inner}</span>
  );
}

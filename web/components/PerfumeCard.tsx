import Link from "next/link";
import { Bottle, BottleStage } from "./Bottle";
import { FamilyBadge } from "./FamilyBadge";
import type { Summary } from "@/lib/types";

type Props = { perfume: Summary; footer?: React.ReactNode; compact?: boolean };

export function PerfumeCard({ perfume, footer, compact }: Props) {
  return (
    <div className="group flex h-full flex-col rounded-2xl border border-line bg-card p-4 transition hover:-translate-y-0.5 hover:shadow-[0_8px_30px_-12px_rgba(60,40,20,0.25)]">
      <Link href={`/perfume/${perfume.id}`} className="flex flex-1 flex-col">
        <BottleStage familyId={perfume.family.id} className="rounded-xl pb-2 pt-4">
          <Bottle id={perfume.id} name={perfume.name} house={perfume.house} familyId={perfume.family.id} size={compact ? 96 : 120} showName />
        </BottleStage>
        <div className="mt-3 flex-1">
          <p className="text-xs uppercase tracking-wider text-faint">{perfume.house}</p>
          <h3 className="font-display text-lg leading-snug group-hover:text-accent">{perfume.name}</h3>
        </div>
        <div className="mt-2">
          <FamilyBadge family={perfume.family} />
        </div>
      </Link>
      {footer && <div className="mt-3 border-t border-line pt-3 text-sm text-muted">{footer}</div>}
    </div>
  );
}

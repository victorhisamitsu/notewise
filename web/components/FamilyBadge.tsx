import { familyColor } from "@/lib/families";
import type { Family } from "@/lib/types";

export function FamilyBadge({ family, size = "sm" }: { family: Family; size?: "sm" | "md" }) {
  if (!family?.name) return null;
  return (
    <span
      className={`inline-flex items-center gap-1.5 rounded-full border border-line bg-card text-muted ${
        size === "md" ? "px-3 py-1 text-sm" : "px-2 py-0.5 text-xs"
      }`}
    >
      <span aria-hidden className="h-2 w-2 rounded-full" style={{ background: familyColor(family.id) }} />
      {family.name}
    </span>
  );
}

import { familyColor, shade } from "@/lib/families";

type Props = {
  id: string;
  name: string;
  house: string;
  familyId?: string;
  concentration?: string;
  size?: number;
  animate?: boolean;
  className?: string;
};

// Deterministic hash so the same perfume always gets the same bottle.
function hash(s: string): number {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return h >>> 0;
}

type Shape = { body: string; clip: string; neckY: number; top: number; bottom: number; left: number; right: number };

// Silhouettes by concentration: slim for toilette, rounded for parfum,
// squat and faceted for extrait / parfum.
function shapeFor(concentration = "", h: number): Shape {
  const c = concentration.toLowerCase();
  if (c.includes("extrait") || c === "parfum" || c.includes("elixir")) {
    const body = "M30 70 L90 70 L100 92 L100 140 Q100 150 90 150 L30 150 Q20 150 20 140 L20 92 Z";
    return { body, clip: body, neckY: 70, top: 70, bottom: 150, left: 20, right: 100 };
  }
  if (c.includes("eau de parfum") || h % 3 === 0) {
    const body = "M40 58 Q24 64 24 84 L24 136 Q24 150 38 150 L82 150 Q96 150 96 136 L96 84 Q96 64 80 58 Z";
    return { body, clip: body, neckY: 58, top: 58, bottom: 150, left: 24, right: 96 };
  }
  const body = "M36 52 Q30 52 30 60 L30 142 Q30 150 38 150 L82 150 Q90 150 90 142 L90 60 Q90 52 84 52 Z";
  return { body, clip: body, neckY: 52, top: 52, bottom: 150, left: 30, right: 90 };
}

export function Bottle({ id, name, house, familyId, concentration, size = 160, animate = false, className }: Props) {
  const h = hash(id);
  const color = familyColor(familyId);
  const shape = shapeFor(concentration, h);
  const uid = `b${h.toString(36)}`;
  const capVariant = h % 3;
  const level = shape.top + (shape.bottom - shape.top) * 0.22; // liquid line
  const span = shape.bottom - level;
  const monogram = (house.trim()[0] ?? "N").toUpperCase();

  return (
    <svg
      viewBox="0 0 120 160"
      width={size}
      height={(size * 160) / 120}
      role="img"
      aria-label={`Illustrated bottle for ${name} by ${house}`}
      className={className}
    >
      <defs>
        <clipPath id={`${uid}-clip`}>
          <path d={shape.clip} />
        </clipPath>
        <linearGradient id={`${uid}-glass`} x1="0" x2="1">
          <stop offset="0" stopColor="#fff" stopOpacity="0.55" />
          <stop offset="0.35" stopColor="#fff" stopOpacity="0.08" />
          <stop offset="1" stopColor="#fff" stopOpacity="0.25" />
        </linearGradient>
      </defs>

      {/* shadow */}
      <ellipse cx="60" cy="153" rx={(shape.right - shape.left) / 2 + 4} ry="4" fill="#000" opacity="0.08" />

      {/* liquid in three layers: top (light), heart, base (deep) */}
      <g clipPath={`url(#${uid}-clip)`}>
        <rect x="0" y={shape.top - 10} width="120" height="120" fill="var(--glass)" />
        <g className={animate ? "fill-up" : undefined}>
          <rect x="0" y={level} width="120" height={span / 3 + 1} fill={shade(color, 0.35)} />
          <rect x="0" y={level + span / 3} width="120" height={span / 3 + 1} fill={color} />
          <rect x="0" y={level + (2 * span) / 3} width="120" height={span / 3 + 2} fill={shade(color, -0.25)} />
          <rect x="0" y={level} width="120" height="2" fill="#fff" opacity="0.45" />
        </g>
        <rect x={shape.left} y={shape.top} width={shape.right - shape.left} height="120" fill={`url(#${uid}-glass)`} />
        <rect x={shape.left + 6} y={shape.top + 8} width="5" height={shape.bottom - shape.top - 20} rx="2.5" fill="#fff" opacity="0.45" />
      </g>

      {/* glass outline */}
      <path d={shape.body} fill="none" stroke={shade(color, -0.45)} strokeOpacity="0.35" strokeWidth="1.5" />

      {/* label with the house monogram */}
      <rect x="44" y={level + span * 0.32} width="32" height="26" rx="3" fill="#fffdf8" opacity="0.92" />
      <text
        x="60"
        y={level + span * 0.32 + 18}
        textAnchor="middle"
        fontFamily="var(--font-display)"
        fontSize="15"
        fill="#3a2f28"
      >
        {monogram}
      </text>

      {/* neck and cap */}
      <rect x="52" y={shape.neckY - 8} width="16" height="9" rx="1.5" fill={shade(color, -0.35)} opacity="0.55" />
      {capVariant === 0 && <rect x="44" y={shape.neckY - 34} width="32" height="27" rx="4" fill="#2f2723" />}
      {capVariant === 1 && <circle cx="60" cy={shape.neckY - 20} r="15" fill="#2f2723" />}
      {capVariant === 2 && (
        <>
          <rect x="38" y={shape.neckY - 22} width="44" height="15" rx="3" fill="#c9a96a" />
          <rect x="38" y={shape.neckY - 22} width="44" height="4" rx="2" fill="#fff" opacity="0.35" />
        </>
      )}
      {capVariant !== 2 && <rect x="48" y={shape.neckY - 30} width="4" height="18" rx="2" fill="#fff" opacity="0.18" />}
    </svg>
  );
}

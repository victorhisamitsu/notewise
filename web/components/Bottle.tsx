import { familyColor, shade } from "@/lib/families";

type Props = {
  id: string;
  name: string;
  house: string;
  familyId?: string;
  concentration?: string;
  size?: number;
  animate?: boolean;
  /** Show the perfume name on the label (readable from ~110px wide). Defaults to size >= 110. */
  showName?: boolean;
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

type Shape = { body: string; neckY: number; top: number; bottom: number; left: number; right: number };

// Silhouettes by concentration: slim for toilette, rounded for parfum,
// squat and faceted for extrait / parfum.
function shapeFor(concentration = "", h: number): Shape {
  const c = concentration.toLowerCase();
  if (c.includes("extrait") || c === "parfum" || c.includes("elixir")) {
    return { body: "M30 70 L90 70 L100 92 L100 140 Q100 150 90 150 L30 150 Q20 150 20 140 L20 92 Z", neckY: 70, top: 70, bottom: 150, left: 20, right: 100 };
  }
  if (c.includes("eau de parfum") || h % 3 === 0) {
    return { body: "M40 58 Q24 64 24 84 L24 136 Q24 150 38 150 L82 150 Q96 150 96 136 L96 84 Q96 64 80 58 Z", neckY: 58, top: 58, bottom: 150, left: 24, right: 96 };
  }
  return { body: "M36 52 Q30 52 30 60 L30 142 Q30 150 38 150 L82 150 Q90 150 90 142 L90 60 Q90 52 84 52 Z", neckY: 52, top: 52, bottom: 150, left: 30, right: 90 };
}

/** Splits a name into at most two lines of ~max characters each. */
function wrap(name: string, max: number): string[] {
  const words = name.split(/\s+/);
  const lines: string[] = [];
  let line = "";
  for (const w of words) {
    if (!line) line = w;
    else if ((line + " " + w).length <= max) line += " " + w;
    else {
      lines.push(line);
      line = w;
    }
  }
  if (line) lines.push(line);
  if (lines.length <= 2) return lines;
  const second = lines.slice(1).join(" ");
  return [lines[0], second.length > max ? second.slice(0, max - 1).trimEnd() + "…" : second];
}

export function Bottle({ id, name, house, familyId, concentration, size = 160, animate = false, showName, className }: Props) {
  const h = hash(id);
  const color = familyColor(familyId);
  const shape = shapeFor(concentration, h);
  const uid = `b${h.toString(36)}`;
  const capVariant = h % 3;
  const width = shape.right - shape.left;
  const level = shape.top + (shape.bottom - shape.top) * 0.2; // liquid line
  const named = showName ?? size >= 110;

  // Label: a cream band across the bottle with the name (or the house initial on small bottles).
  const labelW = Math.min(width - 10, 64);
  const labelX = 60 - labelW / 2;
  const lines = wrap(name, Math.floor(labelW / 4.6));
  const labelH = named ? (lines.length > 1 ? 34 : 26) : 22;
  const labelY = level + (shape.bottom - level) * 0.3;
  const nameSize = lines.some((l) => l.length > Math.floor(labelW / 5.2)) ? 7 : 8;

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
          <path d={shape.body} />
        </clipPath>
        <linearGradient id={`${uid}-liquid`} x1="0" y1="0" x2="0" y2="1">
          <stop offset="0" stopColor={shade(color, 0.3)} />
          <stop offset="0.5" stopColor={color} />
          <stop offset="1" stopColor={shade(color, -0.3)} />
        </linearGradient>
        <linearGradient id={`${uid}-sheen`} x1="0" x2="1">
          <stop offset="0" stopColor="#fff" stopOpacity="0.7" />
          <stop offset="0.18" stopColor="#fff" stopOpacity="0.15" />
          <stop offset="0.55" stopColor="#fff" stopOpacity="0" />
          <stop offset="0.9" stopColor="#fff" stopOpacity="0.2" />
          <stop offset="1" stopColor="#fff" stopOpacity="0.45" />
        </linearGradient>
      </defs>

      {/* soft floor shadow */}
      <ellipse cx="60" cy="153" rx={width / 2 + 6} ry="4.5" fill="#000" opacity="0.1" />

      <g clipPath={`url(#${uid}-clip)`}>
        {/* clear glass: the background shows through */}
        <rect x="0" y="0" width="120" height="160" fill="#fff" fillOpacity="0.28" />
        {/* translucent liquid */}
        <g className={animate ? "fill-up" : undefined}>
          <rect x="0" y={level} width="120" height={shape.bottom - level} fill={`url(#${uid}-liquid)`} fillOpacity="0.78" />
          <rect x="0" y={level} width="120" height="1.6" fill="#fff" opacity="0.6" />
        </g>
        {/* thick glass base */}
        <rect x="0" y={shape.bottom - 9} width="120" height="9" fill="#fff" fillOpacity="0.35" />
        {/* reflections */}
        <rect x="0" y="0" width="120" height="160" fill={`url(#${uid}-sheen)`} />
        <rect x={shape.left + 5} y={shape.top + 7} width="4" height={shape.bottom - shape.top - 22} rx="2" fill="#fff" opacity="0.65" />
        <rect x={shape.right - 8} y={shape.top + 14} width="2" height={shape.bottom - shape.top - 40} rx="1" fill="#fff" opacity="0.5" />
      </g>

      {/* glass edge: a darker outer line and a bright inner line read as thickness */}
      <path d={shape.body} fill="none" stroke={shade(color, -0.55)} strokeOpacity="0.32" strokeWidth="1.4" />
      <path d={shape.body} fill="none" stroke="#fff" strokeOpacity="0.6" strokeWidth="0.6" transform="translate(60 100) scale(0.965) translate(-60 -100)" />

      {/* label */}
      <rect x={labelX} y={labelY} width={labelW} height={labelH} rx="2.5" fill="#fffdf8" fillOpacity="0.94" />
      <rect x={labelX + 2} y={labelY + 2} width={labelW - 4} height={labelH - 4} rx="1.5" fill="none" stroke="#b9a68f" strokeOpacity="0.6" strokeWidth="0.5" />
      {named ? (
        <>
          {lines.map((l, i) => (
            <text
              key={i}
              x="60"
              y={labelY + 11 + i * 9}
              textAnchor="middle"
              fontFamily="var(--font-display)"
              fontSize={nameSize}
              fill="#2f2620"
            >
              {l}
            </text>
          ))}
          <text
            x="60"
            y={labelY + labelH - 5}
            textAnchor="middle"
            fontFamily="var(--font-sans)"
            fontSize="3.6"
            letterSpacing="0.6"
            fill="#7a6a5c"
          >
            {house.toUpperCase().slice(0, 28)}
          </text>
        </>
      ) : (
        <text x="60" y={labelY + 15.5} textAnchor="middle" fontFamily="var(--font-display)" fontSize="13" fill="#2f2620">
          {(house.trim()[0] ?? "N").toUpperCase()}
        </text>
      )}

      {/* neck and cap */}
      <rect x="52" y={shape.neckY - 8} width="16" height="9" rx="1.5" fill="#fff" fillOpacity="0.5" stroke={shade(color, -0.5)} strokeOpacity="0.3" strokeWidth="0.8" />
      {capVariant === 0 && (
        <>
          <rect x="44" y={shape.neckY - 34} width="32" height="27" rx="4" fill="#2f2723" />
          <rect x="47" y={shape.neckY - 31} width="4" height="20" rx="2" fill="#fff" opacity="0.2" />
        </>
      )}
      {capVariant === 1 && (
        <>
          <circle cx="60" cy={shape.neckY - 20} r="14" fill="#fff" fillOpacity="0.35" stroke={shade(color, -0.5)} strokeOpacity="0.35" strokeWidth="1" />
          <circle cx="55" cy={shape.neckY - 25} r="4" fill="#fff" opacity="0.7" />
        </>
      )}
      {capVariant === 2 && (
        <>
          <rect x="38" y={shape.neckY - 22} width="44" height="15" rx="3" fill="#c9a96a" />
          <rect x="38" y={shape.neckY - 22} width="44" height="4" rx="2" fill="#fff" opacity="0.4" />
        </>
      )}
    </svg>
  );
}

/** A soft, family-tinted stage behind a bottle. */
export function BottleStage({ familyId, children, className = "" }: { familyId?: string; children: React.ReactNode; className?: string }) {
  const color = familyColor(familyId);
  return (
    <div
      className={`relative flex items-end justify-center overflow-hidden ${className}`}
      style={{
        background: `radial-gradient(120% 90% at 50% 100%, ${color}55 0%, ${color}22 45%, transparent 75%), linear-gradient(180deg, ${color}10, ${color}26)`,
      }}
    >
      {children}
    </div>
  );
}

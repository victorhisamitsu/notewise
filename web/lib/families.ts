// Visual identity of each olfactory family. These colors tint the
// illustrated bottles and family swatches — decoration, never data encoding.
export const FAMILY_COLORS: Record<string, string> = {
  "sunlit-citrus": "#e8b931",
  "orchard-fruit": "#e06c5a",
  "green-leaf": "#6fa35b",
  "petal-garden": "#e58fb0",
  "white-bloom": "#e2cf9c",
  "powder-veil": "#c3aee0",
  "spice-route": "#c2552d",
  "wood-root": "#8b6a45",
  "smoke-leather": "#5a4a42",
  "amber-resin": "#d08a2e",
  "velvet-sweet": "#a9643c",
  "skin-musk": "#d9bfa8",
  "salt-mineral": "#7fb3c8",
  "oud-court": "#6b3347",
};

export const FALLBACK_COLOR = "#b9a68f";

export function familyColor(id?: string): string {
  return (id && FAMILY_COLORS[id]) || FALLBACK_COLOR;
}

/** Mixes a hex color with white (amount > 0) or black (amount < 0). */
export function shade(hex: string, amount: number): string {
  const n = parseInt(hex.slice(1), 16);
  const mix = (c: number) => {
    const target = amount > 0 ? 255 : 0;
    return Math.round(c + (target - c) * Math.abs(amount));
  };
  const r = mix((n >> 16) & 255);
  const g = mix((n >> 8) & 255);
  const b = mix(n & 255);
  return `#${((1 << 24) | (r << 16) | (g << 8) | b).toString(16).slice(1)}`;
}

export const POPULAR_NOTES = [
  "Vanilla", "Bergamot", "Rose", "Oud", "Sandalwood", "Jasmine", "Patchouli",
  "Amber", "Tonka Bean", "Vetiver", "Lavender", "Musk", "Cedar", "Iris",
  "Pink Pepper", "Leather", "Lemon", "Neroli", "Coffee", "Tuberose",
  "Saffron", "Incense", "Cardamom", "Fig",
];

import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { Bottle } from "./Bottle";

const base = { id: "eau-sauvage-dior", name: "Eau Sauvage", house: "Dior", familyId: "sunlit-citrus" };

describe("Bottle", () => {
  it("is an accessible image named after the perfume", () => {
    render(<Bottle {...base} />);
    expect(screen.getByRole("img", { name: "Illustrated bottle for Eau Sauvage by Dior" })).toBeInTheDocument();
  });

  it("prints the name and house on the label when large enough", () => {
    const { container } = render(<Bottle {...base} size={160} />);
    const text = [...container.querySelectorAll("text")].map((t) => t.textContent).join(" ");
    expect(text).toContain("Eau Sauvage");
    expect(text).toContain("DIOR");
  });

  it("falls back to the house initial on small bottles", () => {
    const { container } = render(<Bottle {...base} size={64} />);
    expect(container.textContent).toBe("D");
  });

  it("wraps long names into at most two lines with an ellipsis", () => {
    const { container } = render(
      <Bottle {...base} id="long" name="Vanilla Royale Sugared Patchouli Extra Long Edition" size={160} />,
    );
    const lines = [...container.querySelectorAll("text")].map((t) => t.textContent ?? "").filter((t) => t !== "DIOR");
    expect(lines.length).toBeLessThanOrEqual(2);
    expect(lines.at(-1)).toMatch(/…$/);
  });

  it("draws the same bottle for the same perfume", () => {
    const a = render(<Bottle {...base} />).container.innerHTML;
    const b = render(<Bottle {...base} />).container.innerHTML;
    expect(a).toBe(b);
  });
});

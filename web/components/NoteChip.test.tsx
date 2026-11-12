import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { NoteChip } from "./NoteChip";

describe("NoteChip", () => {
  it("links to the note page and shows its photo", () => {
    const { container } = render(<NoteChip name="Pink Pepper" />);
    expect(screen.getByRole("link")).toHaveAttribute("href", "/notes/pink-pepper");
    expect(container.querySelector("img")).toHaveAttribute("src", "/notes/pink-pepper.webp");
  });

  it("shows the initial when there is no photo", () => {
    const { container } = render(<NoteChip name="Ambroxan" />);
    expect(container.querySelector("img")).toBeNull();
    expect(screen.getByText("A")).toBeInTheDocument();
  });

  it("can render without a link", () => {
    render(<NoteChip name="Vanilla" link={false} />);
    expect(screen.queryByRole("link")).toBeNull();
    expect(screen.getByText("Vanilla")).toBeInTheDocument();
  });
});

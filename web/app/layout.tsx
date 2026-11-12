import type { Metadata } from "next";
import Link from "next/link";
import "./globals.css";

export const metadata: Metadata = {
  title: { default: "Notewise — discover perfumes by their notes", template: "%s · Notewise" },
  description:
    "Explore perfumes through their note pyramids, find scents by the notes you love, discover the ingredients behind them and get a pick for today's weather.",
};

const nav = [
  { href: "/notes", label: "Note finder" },
  { href: "/ingredients", label: "Ingredients" },
];

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="min-h-screen antialiased">
        <header className="border-b border-line">
          <div className="mx-auto flex max-w-page items-center justify-between gap-4 px-4 py-4 sm:px-6">
            <Link href="/" className="font-display text-2xl tracking-tight">
              Note<span className="italic text-accent">wise</span>
            </Link>
            <nav className="flex gap-1 text-sm">
              {nav.map((n) => (
                <Link key={n.href} href={n.href} className="rounded-full px-3 py-1.5 text-muted transition hover:bg-accent-soft hover:text-ink">
                  {n.label}
                </Link>
              ))}
            </nav>
          </div>
        </header>
        <main className="mx-auto max-w-page px-4 pb-20 sm:px-6">{children}</main>
        <footer className="border-t border-line">
          <div className="mx-auto flex max-w-page flex-col gap-2 px-4 py-8 text-xs text-faint sm:flex-row sm:justify-between sm:px-6">
            <p>
              Perfume data:{" "}
              <a className="underline decoration-line underline-offset-2 hover:text-muted" href="https://www.scento.com/developers/perfume-api" target="_blank" rel="noreferrer">
                Scento Perfume API, scento.com
              </a>
              {" · "}Ingredient photos:{" "}
              <Link className="underline decoration-line underline-offset-2 hover:text-muted" href="/credits">
                Wikimedia Commons
              </Link>
              {" · "}Weather:{" "}
              <a className="underline decoration-line underline-offset-2 hover:text-muted" href="https://open-meteo.com" target="_blank" rel="noreferrer">
                Open-Meteo
              </a>
            </p>
            <p>A portfolio project · Go + Next.js</p>
          </div>
        </footer>
      </body>
    </html>
  );
}

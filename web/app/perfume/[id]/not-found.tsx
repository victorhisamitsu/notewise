import Link from "next/link";
import { SearchBar } from "@/components/SearchBar";

export default function NotFound() {
  return (
    <div className="max-w-xl pt-16">
      <h1 className="font-display text-4xl">We couldn&apos;t find that perfume</h1>
      <p className="mt-3 text-muted">It may not be in the database yet. Try searching by name or house.</p>
      <div className="mt-6">
        <SearchBar />
      </div>
      <Link href="/" className="mt-6 inline-block text-sm text-faint hover:text-accent">
        ← Back home
      </Link>
    </div>
  );
}

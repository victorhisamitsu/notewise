import Link from "next/link";

export default function NotFound() {
  return (
    <div className="max-w-xl pt-16">
      <h1 className="font-display text-4xl">We don&apos;t know that ingredient yet</h1>
      <p className="mt-3 text-muted">Try another spelling, or browse the ingredients we do know.</p>
      <Link href="/ingredients" className="mt-6 inline-block rounded-full bg-ink px-5 py-2.5 text-sm font-medium text-paper hover:opacity-90">
        Browse ingredients
      </Link>
    </div>
  );
}

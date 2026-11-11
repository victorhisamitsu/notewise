"use client";

export default function Error({ reset }: { error: Error; reset: () => void }) {
  return (
    <div className="max-w-xl pt-16">
      <h1 className="font-display text-4xl">The scent trail went cold</h1>
      <p className="mt-3 text-muted">
        The perfume database didn&apos;t answer in time. If the site was idle, the API may be waking up, so try again in a few
        seconds.
      </p>
      <button onClick={reset} className="mt-6 rounded-full bg-ink px-5 py-2.5 text-sm font-medium text-paper hover:opacity-90">
        Try again
      </button>
    </div>
  );
}

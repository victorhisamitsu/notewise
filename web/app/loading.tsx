export default function Loading() {
  return (
    <div className="flex flex-col gap-4 pt-10" aria-label="Loading">
      <div className="h-10 w-2/3 animate-pulse rounded-full bg-card" />
      <div className="h-64 animate-pulse rounded-3xl bg-card" />
    </div>
  );
}

"use client";

import { useState } from "react";
import { clientFetch } from "@/lib/api";
import type { TodayResponse } from "@/lib/types";
import { PerfumeCard } from "./PerfumeCard";

type Weather = { temp: number; humidity: number; code: number; isDay: boolean; place?: string };

// Presets let visitors try the feature without sharing their location.
const PRESETS: { label: string; icon: string; w: Weather }[] = [
  { label: "Hot & humid", icon: "☀️", w: { temp: 31, humidity: 78, code: 1, isDay: true } },
  { label: "Mild day", icon: "⛅", w: { temp: 22, humidity: 50, code: 2, isDay: true } },
  { label: "Rainy", icon: "🌧️", w: { temp: 19, humidity: 90, code: 63, isDay: true } },
  { label: "Cold", icon: "❄️", w: { temp: 7, humidity: 60, code: 3, isDay: true } },
  { label: "Evening out", icon: "🌙", w: { temp: 21, humidity: 55, code: 0, isDay: false } },
];

async function currentWeather(lat: number, lon: number): Promise<Weather> {
  // Called from the browser on purpose: Open-Meteo's free tier is limited
  // per IP, so each visitor uses their own quota.
  const url = `https://api.open-meteo.com/v1/forecast?latitude=${lat.toFixed(2)}&longitude=${lon.toFixed(2)}&current=temperature_2m,relative_humidity_2m,is_day,weather_code`;
  const res = await fetch(url);
  if (!res.ok) throw new Error("weather unavailable");
  const d = await res.json();
  return {
    temp: d.current.temperature_2m,
    humidity: d.current.relative_humidity_2m,
    code: d.current.weather_code,
    isDay: d.current.is_day === 1,
    place: "your location",
  };
}

export function TodayPick() {
  const [state, setState] = useState<"idle" | "locating" | "loading" | "done" | "error">("idle");
  const [weather, setWeather] = useState<Weather | null>(null);
  const [data, setData] = useState<TodayResponse | null>(null);
  const [error, setError] = useState("");

  async function recommend(w: Weather) {
    setWeather(w);
    setState("loading");
    try {
      const q = new URLSearchParams({
        temp: String(Math.round(w.temp)),
        humidity: String(Math.round(w.humidity)),
        code: String(w.code),
        isDay: String(w.isDay),
      });
      setData(await clientFetch<TodayResponse>(`/api/v1/recommendations/today?${q}`));
      setState("done");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Something went wrong");
      setState("error");
    }
  }

  function useLocation() {
    if (!("geolocation" in navigator)) {
      setError("Your browser does not share location. Pick a weather below instead.");
      setState("error");
      return;
    }
    setState("locating");
    navigator.geolocation.getCurrentPosition(
      async (pos) => {
        try {
          await recommend(await currentWeather(pos.coords.latitude, pos.coords.longitude));
        } catch {
          setError("We couldn't read the weather right now. Pick a weather below instead.");
          setState("error");
        }
      },
      () => {
        setError("Location was not shared. Pick a weather below instead.");
        setState("error");
      },
      { timeout: 10000, maximumAge: 600000 },
    );
  }

  return (
    <section aria-labelledby="today" className="rounded-3xl border border-line bg-card p-5 sm:p-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <p className="text-xs uppercase tracking-[0.18em] text-faint">Scent of the day</p>
          <h2 id="today" className="mt-1 font-display text-3xl">What should I wear today?</h2>
          <p className="mt-1 max-w-xl text-sm text-muted">
            Heat makes heavy notes cloying; cold air mutes light ones. Share your local weather or pick one.
          </p>
        </div>
        <button
          type="button"
          onClick={useLocation}
          disabled={state === "locating" || state === "loading"}
          className="shrink-0 rounded-full bg-ink px-5 py-2.5 text-sm font-medium text-paper transition hover:opacity-90 disabled:opacity-50"
        >
          {state === "locating" ? "Reading the sky…" : "📍 Use my weather"}
        </button>
      </div>

      <div className="mt-5 flex flex-wrap gap-2">
        {PRESETS.map((p) => {
          const active = weather?.temp === p.w.temp && weather?.code === p.w.code && weather?.isDay === p.w.isDay && !weather?.place;
          return (
            <button
              key={p.label}
              type="button"
              onClick={() => recommend(p.w)}
              aria-pressed={active}
              className={`rounded-full border px-4 py-2 text-sm transition ${
                active ? "border-accent bg-accent-soft text-accent" : "border-line hover:border-accent"
              }`}
            >
              <span aria-hidden>{p.icon}</span> {p.label}
            </button>
          );
        })}
      </div>

      <div aria-live="polite" className="mt-6">
        {state === "error" && <p className="text-sm text-muted">{error}</p>}
        {state === "loading" && <SkeletonRow />}
        {state === "done" && data && weather && (
          <div className="rise">
            <p className="text-sm text-faint">
              {Math.round(weather.temp)}°C · {Math.round(weather.humidity)}% humidity{weather.place ? ` · ${weather.place}` : ""}
            </p>
            <h3 className="mt-1 font-display text-2xl">{data.headline}</h3>
            <p className="mt-1 max-w-2xl text-muted">{data.reason}</p>
            <ul className="mt-5 grid grid-cols-2 gap-4 md:grid-cols-3">
              {data.picks.map((p) => (
                <li key={p.id}>
                  <PerfumeCard perfume={p} compact footer={<span className="line-clamp-3">{p.why}</span>} />
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
    </section>
  );
}

function SkeletonRow() {
  return (
    <div className="grid grid-cols-2 gap-4 md:grid-cols-3" aria-label="Loading">
      {[0, 1, 2].map((i) => (
        <div key={i} className="h-56 animate-pulse rounded-2xl bg-paper" />
      ))}
    </div>
  );
}

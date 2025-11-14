import type { HeaterProduct } from "../types";

const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8001";
export async function fetchCatalog(): Promise<HeaterProduct[]> {
  const res = await fetch(`${API_BASE}/catalog_heaters`, {
    headers: { Accept: "application/json" },
  });
  if (!res.ok) throw new Error("Failed to fetch catalog");

  const json = await res.json();
  // Берём массив products из data
  return json.data?.products ?? [];
}

export async function fetchHeaterById(id: number | string) {
  const res = await fetch(`${API_BASE}/heater/${id}`, {
    headers: { Accept: "application/json" },
  });
  if (!res.ok) throw new Error("Failed to fetch heater");

  const json = await res.json();
  return json.data?.product ?? json.data ?? json;
}

// src/api.ts
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8001";

export async function fetchCatalog(): Promise<any> {
  // backend: GET /catalog_heaters
  const res = await fetch(`${API_BASE}/catalog_heaters`, { headers: { Accept: "application/json" } });
  // если backend возвращает HTML (rendered templates) — попробуй сделать в бэке JSON endpoint или адаптировать.
  // Предположим он возвращает JSON { products: [...] } — если нет, обработка делается на стороне бэка.
  if (!res.ok) throw new Error("Failed to fetch catalog");
  return res.json();
}

export async function fetchHeaterById(id: number | string) {
  const res = await fetch(`${API_BASE}/heater/${id}`, { headers: { Accept: "application/json" } });
  if (!res.ok) throw new Error("Failed to fetch heater");
  return res.json();
}

export async function fetchApplications(authToken?: string) {
  const headers: any = { Accept: "application/json" };
  if (authToken) headers["Authorization"] = `Bearer ${authToken}`;
  const res = await fetch(`${API_BASE}/heaters_application`, { headers });
  if (!res.ok) throw new Error("Failed to fetch applications");
  return res.json();
}

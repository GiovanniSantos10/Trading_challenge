const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

async function request(path, options) {
  const res = await fetch(`${API_URL}${path}`, options);

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error || `erro ${res.status}`);
  }

  if (res.status === 204) return null;
  return res.json();
}

export function createQuote(mxn) {
  return request("/quotes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ mxn }),
  });
}

export function confirmQuote(id) {
  return request(`/quotes/${id}/confirm`, { method: "POST" });
}

export function getHistory() {
  return request("/history");
}

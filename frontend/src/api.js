// Camada de acesso à API — só fetch nativo, sem axios (regra do projeto).
// A URL base vem de uma variável de ambiente (VITE_API_URL), nunca hardcoded
// no meio do código: assim o front funciona em qualquer porta/host da API
// sem precisar editar JS, só o .env.
const API_URL = import.meta.env.VITE_API_URL || "http://localhost:8080";

// request centraliza o tratamento de erro: toda resposta HTTP >= 400 vira
// um Error com a mensagem que a API já manda em {"error": "..."} , em vez de
// cada chamada repetir esse "if (!res.ok) ...".
async function request(path, options) {
  const res = await fetch(`${API_URL}${path}`, options);

  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error || `erro ${res.status}`);
  }

  if (res.status === 204) return null;
  return res.json();
}

// createQuote pede uma cotação nova: POST /quotes {"mxn": number}.
// Retorna {id, mxn, brl, created_at, expires_at, status}.
export function createQuote(mxn) {
  return request("/quotes", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ mxn }),
  });
}

// confirmQuote confirma a cotação pendente pelo id, antes dela expirar (10s).
export function confirmQuote(id) {
  return request(`/quotes/${id}/confirm`, { method: "POST" });
}

// getHistory lista só as cotações confirmadas (regra da API).
export function getHistory() {
  return request("/history");
}

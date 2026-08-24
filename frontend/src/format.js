// Intl.NumberFormat é nativo do JS (sem precisar de lib) — formata valores
// monetários de acordo com a moeda, sem cada componente reimplementar isso.
const brl = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const mxn = new Intl.NumberFormat("es-MX", { style: "currency", currency: "MXN" });

export function formatBRL(value) {
  return brl.format(value);
}

export function formatMXN(value) {
  return mxn.format(value);
}

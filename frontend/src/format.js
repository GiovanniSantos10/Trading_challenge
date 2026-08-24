const brl = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const mxn = new Intl.NumberFormat("es-MX", { style: "currency", currency: "MXN" });
const rate = new Intl.NumberFormat("pt-BR", { minimumFractionDigits: 4, maximumFractionDigits: 4 });

export function formatBRL(value) {
  return brl.format(value);
}

export function formatMXN(value) {
  return mxn.format(value);
}

export function formatRate(value) {
  return rate.format(value);
}

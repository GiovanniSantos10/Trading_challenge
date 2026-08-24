// Intl.NumberFormat é nativo do JS (sem precisar de lib) — formata valores
// monetários de acordo com a moeda, sem cada componente reimplementar isso.
const brl = new Intl.NumberFormat("pt-BR", { style: "currency", currency: "BRL" });
const mxn = new Intl.NumberFormat("es-MX", { style: "currency", currency: "MXN" });

// Preço por unidade (BRL/MXN) não é um valor monetário "fechado" — é uma
// taxa, geralmente menor que 1. Formatar como moeda arredondaria pra 2 casas
// (perderia precisão real do câmbio), então usa mais casas decimais, sem
// símbolo de moeda embutido (o BRL/MXN é adicionado onde é exibido).
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

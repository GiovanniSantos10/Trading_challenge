import { useEffect, useState } from "react";
import { formatBRL, formatMXN } from "./format";

// QuoteResult é dona do contador regressivo: recalcula secondsLeft a partir
// de quote.expires_at (o servidor é a fonte da verdade do prazo, o front só
// exibe) em vez de contar 10 "na unha" — assim um delay de rede não desalinha
// o contador com o que a API realmente vai aceitar.
//
// "now" (não "secondsLeft") é o único estado guardado, e só é atualizado
// dentro do callback do setInterval — nunca direto no corpo do efeito. Isso
// evita o padrão "setState síncrono dentro de um efeito" (o oxlint reclama
// dele com razão: dispara um re-render em cascata extra logo após o efeito
// rodar); secondsLeft/expired são derivados durante o render a partir de
// "now" e da prop quote, não guardados em estado próprio.
function QuoteResult({ quote, onConfirm, onExpire, confirming }) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!quote) return;
    const intervalId = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(intervalId);
  }, [quote]);

  const secondsLeft = quote
    ? Math.max(0, Math.round((new Date(quote.expires_at).getTime() - now) / 1000))
    : 0;
  const expired = Boolean(quote) && secondsLeft <= 0;

  // Avisar o pai (que decide o que fazer — aqui, limpar a cotação) é o
  // "efeito colateral" de verdade; ele acontece fora do render, não dentro
  // dele, mesmo a condição sendo derivada durante o render.
  useEffect(() => {
    if (expired) onExpire();
  }, [expired, onExpire]);

  if (!quote) return null;

  return (
    <div className="quote-result">
      <p className="quote-amount">
        {formatMXN(quote.mxn)} = <strong>{formatBRL(quote.brl)}</strong>
      </p>
      <p className={expired ? "countdown countdown-expired" : "countdown"}>
        {expired ? "Cotação expirada" : `Expira em ${secondsLeft}s`}
      </p>
      <button onClick={() => onConfirm(quote.id)} disabled={expired || confirming}>
        {confirming ? "Confirmando..." : "Confirmar cotação"}
      </button>
    </div>
  );
}

export default QuoteResult;

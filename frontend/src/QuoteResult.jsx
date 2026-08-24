import { useEffect, useState } from "react";
import { formatBRL } from "./format";

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

  // Cor do badge muda conforme o prazo aperta — pista visual de urgência,
  // igual ao que corretoras cripto usam pra cotação com validade curta.
  const urgency = secondsLeft <= 3 ? "urgent" : secondsLeft <= 6 ? "warning" : "ok";

  return (
    <div className="swap-row">
      <div className="row-label-line">
        <span className="row-label">Você paga</span>
        {quote && (
          <span className={`countdown-pill countdown-${expired ? "urgent" : urgency}`}>
            {expired ? "expirada" : `${secondsLeft}s`}
          </span>
        )}
      </div>
      <div className="swap-row-main">
        <span className={quote ? "swap-amount" : "swap-amount placeholder"}>
          {quote ? formatBRL(quote.brl) : "0,00"}
        </span>
        <span className="token-pill">BRL</span>
      </div>

      {quote && (
        <button className="cta" onClick={() => onConfirm(quote.id)} disabled={expired || confirming}>
          {confirming ? "Confirmando..." : "Confirmar cotação"}
        </button>
      )}
    </div>
  );
}

export default QuoteResult;

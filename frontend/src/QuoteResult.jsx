import { useEffect, useState } from "react";
import { formatBRL, formatRate } from "./format";

// QuoteResult é dona do contador regressivo: recalcula secondsLeft a partir
// de quote.expires_at (o servidor é a fonte da verdade do prazo, o front só
// exibe) em vez de contar 10 "na unha" — assim um delay de rede não desalinha
// o contador com o que a API realmente vai aceitar.
//
// "now" (não "secondsLeft") é o único estado guardado; secondsLeft/expired
// são derivados durante o render a partir dele e da prop quote, não
// guardados em estado próprio.
//
// O setNow(Date.now()) logo no início do efeito é síncrono de propósito —
// e sim, isso é o padrão que o oxlint (react(set-state-in-effect)) normalmente
// reclama. Uma versão anterior tirou essa chamada pra calar o aviso, achando
// que "now" continuaria fresco sozinho (o componente fica montado o tempo
// todo, então o intervalo de uma cotação anterior mantinha o relógio girando).
// Só que enquanto NENHUMA cotação existe, o efeito sai cedo (if (!quote)
// return) e o intervalo nem liga — "now" fica parado no valor do mount. Se o
// usuário demora pra pedir a primeira cotação, esse atraso vaza pro cálculo:
// o contador nasce inflado e "salta" pro valor certo 1s depois, quando o
// intervalo finalmente liga. Por isso o resync sempre que "quote" muda é
// necessário pra corretude, não só uma otimização — o aviso do lint é um
// falso positivo aqui.
function QuoteResult({ quote, onConfirm, onExpire, confirming }) {
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (!quote) return;
    setNow(Date.now());
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

      {quote && <p className="rate-note">{formatRate(quote.brl / quote.mxn)} BRL/MXN</p>}

      {quote && (
        <button className="cta" onClick={() => onConfirm(quote.id)} disabled={expired || confirming}>
          {confirming ? "Confirmando..." : "Confirmar cotação"}
        </button>
      )}
    </div>
  );
}

export default QuoteResult;

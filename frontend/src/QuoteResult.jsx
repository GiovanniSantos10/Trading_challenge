import { useEffect, useState } from "react";
import { formatBRL, formatRate } from "./format";

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

  useEffect(() => {
    if (expired) onExpire();
  }, [expired, onExpire]);

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

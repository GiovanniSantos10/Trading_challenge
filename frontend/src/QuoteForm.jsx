import { useState } from "react";

// QuoteForm não guarda o resultado da cotação — só coleta o valor em MXN e
// avisa o componente pai (onQuote) quando o usuário confirma o envio. Quem
// decide o que fazer com a cotação é o App, não este componente.
//
// "locked" (true enquanto existe uma cotação pendente) esconde o botão
// "Cotar" e desabilita o input — igual ao passo de "review" de uma corretora,
// o valor fica travado enquanto a cotação atual está valendo.
function QuoteForm({ onQuote, loading, locked }) {
  const [mxn, setMxn] = useState("");

  function handleSubmit(event) {
    event.preventDefault();

    const value = Number(mxn);
    if (!value || value <= 0) return;

    onQuote(value);
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="swap-row">
        <label htmlFor="mxn-input" className="row-label">
          Valor a receber
        </label>
        <div className="swap-row-main">
          <input
            id="mxn-input"
            type="number"
            min="0"
            step="0.01"
            placeholder="0,00"
            value={mxn}
            onChange={(event) => setMxn(event.target.value)}
            disabled={loading || locked}
          />
          <span className="token-pill">MXN</span>
        </div>
      </div>

      {!locked && (
        <button type="submit" className="cta" disabled={loading || !mxn}>
          {loading ? "Cotando..." : "Cotar"}
        </button>
      )}
    </form>
  );
}

export default QuoteForm;

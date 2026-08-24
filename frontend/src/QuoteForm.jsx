import { useState } from "react";

// QuoteForm não guarda o resultado da cotação — só coleta o valor em MXN e
// avisa o componente pai (onQuote) quando o usuário confirma o envio. Quem
// decide o que fazer com a cotação é o App, não este componente.
function QuoteForm({ onQuote, loading }) {
  const [mxn, setMxn] = useState("");

  function handleSubmit(event) {
    event.preventDefault();

    const value = Number(mxn);
    if (!value || value <= 0) return;

    onQuote(value);
  }

  return (
    <form onSubmit={handleSubmit} className="quote-form">
      <label htmlFor="mxn">Valor em MXN</label>
      <div className="quote-form-row">
        <input
          id="mxn"
          type="number"
          min="0"
          step="0.01"
          placeholder="1000"
          value={mxn}
          onChange={(event) => setMxn(event.target.value)}
          disabled={loading}
        />
        <button type="submit" disabled={loading}>
          {loading ? "Cotando..." : "Cotar"}
        </button>
      </div>
    </form>
  );
}

export default QuoteForm;

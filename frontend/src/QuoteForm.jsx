import { useState } from "react";

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

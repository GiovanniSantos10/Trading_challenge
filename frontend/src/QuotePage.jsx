import QuoteForm from "./QuoteForm";
import QuoteResult from "./QuoteResult";

function QuotePage({ quote, loading, confirming, error, onQuote, onConfirm, onExpire }) {
  return (
    <div className="swap-card">
      <div className="swap-card-header">
        <h1>Cotador</h1>
        <span className="pair-pill">MXN → BRL</span>
      </div>

      <QuoteForm onQuote={onQuote} loading={loading} locked={Boolean(quote)} />

      <div className="swap-divider" aria-hidden="true">
        <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2.5">
          <path d="M6 9l6 6 6-6" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </div>

      <QuoteResult quote={quote} onConfirm={onConfirm} onExpire={onExpire} confirming={confirming} />

      {error && <p className="error">{error}</p>}
    </div>
  );
}

export default QuotePage;

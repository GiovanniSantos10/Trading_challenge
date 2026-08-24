import { formatMXN, formatRate } from "./format";

function History({ quotes }) {
  if (quotes.length === 0) {
    return <p className="history-empty">Nenhuma cotação confirmada ainda.</p>;
  }

  return (
    <ul className="history-list">
      {quotes.map((q) => (
        <li key={q.id} className="history-item">
          <div className="history-item-pair">
            <span className="history-item-icon" aria-hidden="true">
              →
            </span>
            <div>
              <div className="history-item-amounts">{formatMXN(q.mxn)}</div>
              <div className="history-item-time">{new Date(q.created_at).toLocaleString("pt-BR")}</div>
            </div>
          </div>
          <span className="history-item-rate">{formatRate(q.brl / q.mxn)} BRL/MXN</span>
        </li>
      ))}
    </ul>
  );
}

export default History;

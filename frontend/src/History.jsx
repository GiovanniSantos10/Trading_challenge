import { formatBRL, formatMXN } from "./format";

// History só recebe a lista e desenha — não busca dados sozinha. Quem chama
// GET /history e decide quando atualizar é o App (é ele que sabe quando uma
// confirmação acabou de acontecer).
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
              <div className="history-item-amounts">
                {formatMXN(q.mxn)} <span className="history-item-arrow">→</span> {formatBRL(q.brl)}
              </div>
              <div className="history-item-time">{new Date(q.created_at).toLocaleString("pt-BR")}</div>
            </div>
          </div>
          <span className="history-item-status">confirmada</span>
        </li>
      ))}
    </ul>
  );
}

export default History;

import { formatMXN, formatRate } from "./format";

// History só recebe a lista e desenha — não busca dados sozinha. Quem chama
// GET /history e decide quando atualizar é o App (é ele que sabe quando uma
// confirmação acabou de acontecer).
//
// O enunciado pede quantidade + preço por unidade + data/hora. "Preço por
// unidade" (a taxa de câmbio usada, BRL por MXN) não vem pronto da API — é
// derivado aqui (brl / mxn) a partir dos dois valores que já existem, sem
// precisar mudar o backend só pra isso.
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

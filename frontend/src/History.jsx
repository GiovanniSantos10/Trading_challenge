import { formatBRL, formatMXN } from "./format";

// History só recebe a lista e desenha — não busca dados sozinha. Quem chama
// GET /history e decide quando atualizar é o App (é ele que sabe quando uma
// confirmação acabou de acontecer).
function History({ quotes }) {
  if (quotes.length === 0) {
    return <p className="history-empty">Nenhuma cotação confirmada ainda.</p>;
  }

  return (
    <table className="history-table">
      <thead>
        <tr>
          <th>MXN</th>
          <th>BRL</th>
          <th>Confirmada em</th>
        </tr>
      </thead>
      <tbody>
        {quotes.map((q) => (
          <tr key={q.id}>
            <td>{formatMXN(q.mxn)}</td>
            <td>{formatBRL(q.brl)}</td>
            <td>{new Date(q.created_at).toLocaleString("pt-BR")}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

export default History;

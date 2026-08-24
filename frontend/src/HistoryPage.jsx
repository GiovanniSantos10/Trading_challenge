import History from "./History";

// HistoryPage é a segunda página pedida no enunciado: lista as cotações
// confirmadas. Sem filtro, sem paginação — só recebe a lista pronta do App.
function HistoryPage({ quotes }) {
  return (
    <div className="swap-card">
      <div className="swap-card-header">
        <h1>Histórico</h1>
        <span className="pair-pill">confirmadas</span>
      </div>
      <History quotes={quotes} />
    </div>
  );
}

export default HistoryPage;

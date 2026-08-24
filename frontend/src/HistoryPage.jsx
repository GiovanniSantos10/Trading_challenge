import History from "./History";

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

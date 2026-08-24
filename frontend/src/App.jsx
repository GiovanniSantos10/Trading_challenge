import { useCallback, useEffect, useState } from "react";
import { createQuote, confirmQuote, getHistory } from "./api";
import QuoteForm from "./QuoteForm";
import QuoteResult from "./QuoteResult";
import History from "./History";
import "./App.css";

function App() {
  // Cada pedaço de UI vem de uma variável de estado — nada fixo no JSX.
  const [quote, setQuote] = useState(null); // cotação pendente atual, ou null
  const [history, setHistory] = useState([]); // cotações confirmadas
  const [loading, setLoading] = useState(false); // criando cotação
  const [confirming, setConfirming] = useState(false); // confirmando cotação
  const [error, setError] = useState(null);

  // useCallback: a função precisa ter identidade estável entre renders,
  // porque QuoteResult a usa como dependência do useEffect do contador —
  // sem isso, o contador reiniciaria o intervalo a cada render do App.
  const refreshHistory = useCallback(async () => {
    try {
      const list = await getHistory();
      setHistory(list);
    } catch {
      // histórico é secundário: se falhar, não interrompe o fluxo principal.
    }
  }, []);

  // Carrega o histórico uma vez, quando o componente monta.
  useEffect(() => {
    refreshHistory();
  }, [refreshHistory]);

  async function handleQuote(mxn) {
    setLoading(true);
    setError(null);
    try {
      const result = await createQuote(mxn);
      setQuote(result);
    } catch (err) {
      setError(err.message);
      setQuote(null);
    } finally {
      setLoading(false);
    }
  }

  async function handleConfirm(id) {
    setConfirming(true);
    setError(null);
    try {
      await confirmQuote(id);
      setQuote(null);
      await refreshHistory();
    } catch (err) {
      setError(err.message);
    } finally {
      setConfirming(false);
    }
  }

  const handleExpire = useCallback(() => {
    setQuote(null);
  }, []);

  return (
    <main className="app">
      <div className="swap-card">
        <div className="swap-card-header">
          <h1>Cotador</h1>
          <span className="pair-pill">MXN → BRL</span>
        </div>

        <QuoteForm onQuote={handleQuote} loading={loading} locked={Boolean(quote)} />

        {/* Só indica que o fluxo desce de uma linha pra outra — não é um
            botão, cotação aqui é sempre nessa única direção (sem alternar). */}
        <div className="swap-divider" aria-hidden="true">
          <svg viewBox="0 0 24 24" width="14" height="14" fill="none" stroke="currentColor" strokeWidth="2.5">
            <path d="M6 9l6 6 6-6" strokeLinecap="round" strokeLinejoin="round" />
          </svg>
        </div>

        <QuoteResult
          quote={quote}
          onConfirm={handleConfirm}
          onExpire={handleExpire}
          confirming={confirming}
        />

        {error && <p className="error">{error}</p>}
      </div>

      <section className="history-section">
        <h2>Atividade recente</h2>
        <History quotes={history} />
      </section>
    </main>
  );
}

export default App;

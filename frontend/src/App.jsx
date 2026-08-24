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
      <h1>Cotador BRL / MXN</h1>
      <p className="subtitle">Informe o valor em MXN e receba a cotação em BRL, via USDT.</p>

      <QuoteForm onQuote={handleQuote} loading={loading} />

      {error && <p className="error">{error}</p>}

      <QuoteResult
        quote={quote}
        onConfirm={handleConfirm}
        onExpire={handleExpire}
        confirming={confirming}
      />

      <section className="history-section">
        <h2>Histórico</h2>
        <History quotes={history} />
      </section>
    </main>
  );
}

export default App;

import { useCallback, useEffect, useState } from "react";
import { createQuote, confirmQuote, getHistory } from "./api";
import { useRoute } from "./useRoute";
import Nav from "./Nav";
import QuotePage from "./QuotePage";
import HistoryPage from "./HistoryPage";
import "./App.css";

function App() {
  // Rota atual ("/" ou "/historico") — as duas páginas pedidas no enunciado.
  const [path, navigate] = useRoute();

  // Cada pedaço de UI vem de uma variável de estado — nada fixo no JSX. Tudo
  // fica aqui no App (não em cada página) porque trocar de rota não
  // desmonta o App — o estado sobrevive à navegação entre as duas páginas.
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
      <Nav path={path} onNavigate={navigate} />

      {path === "/historico" ? (
        <HistoryPage quotes={history} />
      ) : (
        <QuotePage
          quote={quote}
          loading={loading}
          confirming={confirming}
          error={error}
          onQuote={handleQuote}
          onConfirm={handleConfirm}
          onExpire={handleExpire}
        />
      )}
    </main>
  );
}

export default App;

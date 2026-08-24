import { useCallback, useEffect, useState } from "react";
import { createQuote, confirmQuote, getHistory } from "./api";
import { useRoute } from "./useRoute";
import Nav from "./Nav";
import QuotePage from "./QuotePage";
import HistoryPage from "./HistoryPage";
import "./App.css";

function App() {
  const [path, navigate] = useRoute();

  const [quote, setQuote] = useState(null);
  const [history, setHistory] = useState([]);
  const [loading, setLoading] = useState(false);
  const [confirming, setConfirming] = useState(false);
  const [error, setError] = useState(null);

  const refreshHistory = useCallback(async () => {
    try {
      const list = await getHistory();
      setHistory(list);
    } catch {}
  }, []);

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

import { useCallback, useEffect, useState } from "react";

// Roteador mínimo, sem lib nova (regra do projeto): a URL de verdade
// (window.location) é o "sistema externo" com o qual este hook sincroniza —
// exatamente o caso de uso que a documentação do React usa como exemplo
// legítimo de useEffect. "popstate" é o evento nativo do browser disparado
// quando o usuário usa os botões voltar/avançar.
export function useRoute() {
  const [path, setPath] = useState(window.location.pathname);

  useEffect(() => {
    function onPopState() {
      setPath(window.location.pathname);
    }
    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  // pushState muda a URL sem recarregar a página; setPath (fora do evento
  // popstate) é o que faz o React re-renderizar com a nova rota.
  const navigate = useCallback((to) => {
    if (to !== window.location.pathname) {
      window.history.pushState({}, "", to);
    }
    setPath(to);
  }, []);

  return [path, navigate];
}

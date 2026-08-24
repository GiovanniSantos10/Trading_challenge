import { useCallback, useEffect, useState } from "react";

export function useRoute() {
  const [path, setPath] = useState(window.location.pathname);

  useEffect(() => {
    function onPopState() {
      setPath(window.location.pathname);
    }
    window.addEventListener("popstate", onPopState);
    return () => window.removeEventListener("popstate", onPopState);
  }, []);

  const navigate = useCallback((to) => {
    if (to !== window.location.pathname) {
      window.history.pushState({}, "", to);
    }
    setPath(to);
  }, []);

  return [path, navigate];
}

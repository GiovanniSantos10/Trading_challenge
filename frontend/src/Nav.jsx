// Nav é "burro": só recebe a rota atual e avisa o pai quando o usuário quer
// navegar. <a href> de verdade (com preventDefault) em vez de <button> —
// mantém foco/teclado/"abrir em nova aba" funcionando como um link normal.
function Nav({ path, onNavigate }) {
  function linkProps(to) {
    return {
      href: to,
      className: path === to ? "nav-link nav-link-active" : "nav-link",
      onClick: (event) => {
        event.preventDefault();
        onNavigate(to);
      },
    };
  }

  return (
    <nav className="nav">
      <a {...linkProps("/")}>Cotação</a>
      <a {...linkProps("/historico")}>Histórico</a>
    </nav>
  );
}

export default Nav;

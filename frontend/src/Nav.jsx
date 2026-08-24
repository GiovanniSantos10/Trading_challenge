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

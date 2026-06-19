import { Link, useNavigate } from "react-router-dom";

function Navbar({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    navigate("/");
  };

  return (
    <nav className="navbar">
      <Link to="/" className="navbar__logo">
        TicketHub
      </Link>

      <div className="navbar__links">
        <Link to="/" className="navbar__link">Eventos</Link>

        {user ? (
          <>
            <Link to="/tickets" className="navbar__link">Mis Entradas</Link>
            {user.rol === "admin" && (
              <Link to="/admin" className="navbar__link navbar__link--admin">Admin</Link>
            )}
            <div className="navbar__user">{user.email}</div>
            <button className="btn btn--ghost" onClick={handleLogout}>Salir</button>
          </>
        ) : (
          <>
            <Link to="/login" className="navbar__link">Iniciar sesión</Link>
            <Link to="/register" className="btn btn--primary btn--sm">Registrarse</Link>
          </>
        )}
      </div>
    </nav>
  );
}

export default Navbar;

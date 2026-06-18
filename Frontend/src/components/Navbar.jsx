import { Link, useNavigate } from "react-router-dom";

function Navbar({ user, onLogout }) {
  const navigate = useNavigate();

  const handleLogout = () => {
    onLogout();
    navigate("/login");
  };

  return (
    <nav>
      <Link to="/">Inicio</Link>
      {user ? (
        <>
          {" | "}
          <Link to="/tickets">Mis Entradas</Link>
          {" | "}
          <span>{user.email}</span>
          {" | "}
          <button onClick={handleLogout}>Cerrar sesión</button>
        </>
      ) : (
        <>
          {" | "}
          <Link to="/login">Login</Link>
          {" | "}
          <Link to="/register">Registro</Link>
        </>
      )}
    </nav>
  );
}

export default Navbar;

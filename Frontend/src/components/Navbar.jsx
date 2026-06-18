import { Link } from "react-router-dom";

function Navbar() {
  return (
    <nav>
      <Link to="/">Inicio</Link> |{" "}
      <Link to="/tickets">Mis Entradas</Link> |{" "}
      <Link to="/login">Login</Link> |{" "}
      <Link to="/register">Registro</Link>
    </nav>
  );
}

export default Navbar;
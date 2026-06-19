import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { register } from "../services/authService";

function Register() {
  const [form, setForm]       = useState({ nombre: "", apellido: "", email: "", password: "", dni: "" });
  const [error, setError]     = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      await register(form);
      navigate("/login");
    } catch (err) {
      setError(err.response?.data?.error || "Error al registrarse. Intentá de nuevo.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card" style={{ maxWidth: 460 }}>
        <div className="auth-card__logo">TH</div>
        <h1 className="auth-card__title">Crear cuenta</h1>
        <p className="auth-card__subtitle">Completá tus datos para registrarte</p>

        {error && <div className="alert alert--error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="nombre">Nombre</label>
              <input
                id="nombre"
                className="form-input"
                name="nombre"
                placeholder="Juan"
                value={form.nombre}
                onChange={handleChange}
                autoComplete="given-name"
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="apellido">Apellido</label>
              <input
                id="apellido"
                className="form-input"
                name="apellido"
                placeholder="Pérez"
                value={form.apellido}
                onChange={handleChange}
                autoComplete="family-name"
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="email">Email</label>
            <input
              id="email"
              className="form-input"
              name="email"
              type="email"
              placeholder="tu@email.com"
              value={form.email}
              onChange={handleChange}
              autoComplete="email"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="password">Contraseña</label>
            <input
              id="password"
              className="form-input"
              name="password"
              type="password"
              placeholder="••••••••"
              value={form.password}
              onChange={handleChange}
              autoComplete="new-password"
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="dni">DNI</label>
            <input
              id="dni"
              className="form-input"
              name="dni"
              placeholder="12345678"
              value={form.dni}
              onChange={handleChange}
              required
            />
          </div>

          <button
            type="submit"
            className="btn btn--primary btn--full btn--lg"
            style={{ marginTop: 8 }}
            disabled={loading}
          >
            {loading ? "Registrando..." : "Crear cuenta"}
          </button>
        </form>

        <p style={{ textAlign: "center", marginTop: 24, fontSize: 14, color: "var(--text-muted)" }}>
          ¿Ya tenés cuenta?{" "}
          <Link to="/login">Iniciá sesión</Link>
        </p>
      </div>
    </div>
  );
}

export default Register;

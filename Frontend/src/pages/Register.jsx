import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { register } from "../services/authService";

export default function Register() {
  const [form, setForm] = useState({ nombre: "", apellido: "", email: "", password: "", dni: "" });
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const handleChange = (e) => setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError(""); setLoading(true);
    try {
      await register(form);
      navigate("/login");
    } catch (err) {
      setError(err.response?.data?.error || "Error al registrarse.");
    } finally { setLoading(false); }
  };

  return (
    <div className="auth-page">
      <div className="auth-card" style={{ maxWidth: 460 }}>
        <h1 className="auth-card__title">Crear cuenta</h1>

        {error && <div className="alert alert--error">{error}</div>}

        <form onSubmit={handleSubmit}>
          <div className="form-row">
            <div className="form-group">
              <label className="form-label">Nombre</label>
              <input className="form-input" name="nombre" placeholder="Juan" value={form.nombre} onChange={handleChange} required />
            </div>
            <div className="form-group">
              <label className="form-label">Apellido</label>
              <input className="form-input" name="apellido" placeholder="Pérez" value={form.apellido} onChange={handleChange} required />
            </div>
          </div>
          <div className="form-group">
            <label className="form-label">Email</label>
            <input className="form-input" name="email" type="email" placeholder="tu@email.com" value={form.email} onChange={handleChange} required />
          </div>
          <div className="form-group">
            <label className="form-label">Contraseña</label>
            <input className="form-input" name="password" type="password" placeholder="Mínimo 6 caracteres" value={form.password} onChange={handleChange} required />
          </div>
          <div className="form-group">
            <label className="form-label">DNI</label>
            <input className="form-input" name="dni" placeholder="12345678" value={form.dni} onChange={handleChange} required />
          </div>
          <button type="submit" className="btn btn--primary btn--full btn--lg" disabled={loading} style={{ marginTop: 8 }}>
            {loading ? "Registrando..." : "Crear cuenta"}
          </button>
        </form>

        <p style={{ textAlign: "center", marginTop: 24, fontSize: 14 }}>
          ¿Ya tenés cuenta? <Link to="/login">Iniciá sesión</Link>
        </p>
      </div>
    </div>
  );
}

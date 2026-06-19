import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getEvents, createEvent, deleteEvent } from "../services/eventService";

const EMPTY_FORM = {
  titulo: "",
  descripcion: "",
  categoria: "Recitales",
  fecha: "",
  hora: "",
  lugar: "",
  precio: "",
  cupo_maximo: "",
};

function AdminPanel({ user }) {
  const navigate = useNavigate();
  const [events, setEvents]     = useState([]);
  const [form, setForm]         = useState(EMPTY_FORM);
  const [loading, setLoading]   = useState(false);
  const [fetching, setFetching] = useState(true);
  const [success, setSuccess]   = useState("");
  const [error, setError]       = useState("");

  useEffect(() => {
    if (!user || user.rol !== "admin") {
      navigate("/");
      return;
    }
    cargarEventos();
  }, [user]);

  const cargarEventos = () => {
    setFetching(true);
    getEvents()
      .then((res) => setEvents(res.data))
      .catch(() => setError("No se pudieron cargar los eventos"))
      .finally(() => setFetching(false));
  };

  const handleChange = (e) =>
    setForm({ ...form, [e.target.name]: e.target.value });

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setSuccess("");
    setLoading(true);
    try {
      await createEvent({
        ...form,
        precio: parseFloat(form.precio),
        cupo_maximo: parseInt(form.cupo_maximo),
        fecha: new Date(`${form.fecha}T${form.hora || "00:00"}`).toISOString(),
      });
      setSuccess(`Evento "${form.titulo}" creado correctamente.`);
      setForm(EMPTY_FORM);
      cargarEventos();
    } catch (err) {
      setError(err.response?.data?.error || "Error al crear el evento");
    } finally {
      setLoading(false);
    }
  };

  const handleEliminar = async (evento) => {
    if (!window.confirm(`¿Eliminár "${evento.titulo}"? Esta acción no se puede deshacer.`)) return;
    setError("");
    setSuccess("");
    try {
      await deleteEvent(evento.id);
      setSuccess(`Evento "${evento.titulo}" eliminado.`);
      cargarEventos();
    } catch (err) {
      setError(err.response?.data?.error || "Error al eliminar el evento");
    }
  };

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1 className="page-title">Panel de Administración</h1>
        <p className="page-subtitle">Creá y gestioná los eventos del sistema</p>
      </div>

      {success && <div className="alert alert--success">{success}</div>}
      {error   && <div className="alert alert--error">{error}</div>}

      {/* Formulario de nuevo evento */}
      <div className="admin-section">
        <h2 className="admin-section__title">Nuevo Evento</h2>
        <form onSubmit={handleSubmit} className="admin-form">
          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="titulo">Título</label>
              <input
                id="titulo"
                className="form-input"
                name="titulo"
                placeholder="Ej: Coldplay en Córdoba"
                value={form.titulo}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="categoria">Categoría</label>
              <select
                id="categoria"
                className="form-input"
                name="categoria"
                value={form.categoria}
                onChange={handleChange}
                required
              >
                <option value="Recitales">Recitales</option>
                <option value="Teatro">Teatro</option>
                <option value="Deportes">Deportes</option>
              </select>
            </div>
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="descripcion">Descripción</label>
            <textarea
              id="descripcion"
              className="form-input form-textarea"
              name="descripcion"
              placeholder="Descripción del evento..."
              value={form.descripcion}
              onChange={handleChange}
              rows={3}
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="fecha">Fecha</label>
              <input
                id="fecha"
                className="form-input"
                name="fecha"
                type="date"
                value={form.fecha}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="hora">Horario</label>
              <input
                id="hora"
                className="form-input"
                name="hora"
                type="time"
                value={form.hora}
                onChange={handleChange}
                required
              />
            </div>
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="lugar">Lugar</label>
            <input
              id="lugar"
              className="form-input"
              name="lugar"
              placeholder="Ej: Estadio Mario Kempes"
              value={form.lugar}
              onChange={handleChange}
              required
            />
          </div>

          <div className="form-row">
            <div className="form-group">
              <label className="form-label" htmlFor="precio">Precio ($)</label>
              <input
                id="precio"
                className="form-input"
                name="precio"
                type="number"
                min="1"
                step="0.01"
                placeholder="5000"
                value={form.precio}
                onChange={handleChange}
                required
              />
            </div>
            <div className="form-group">
              <label className="form-label" htmlFor="cupo_maximo">Cupo máximo</label>
              <input
                id="cupo_maximo"
                className="form-input"
                name="cupo_maximo"
                type="number"
                min="1"
                placeholder="200"
                value={form.cupo_maximo}
                onChange={handleChange}
                required
              />
            </div>
          </div>

          <button
            type="submit"
            className="btn btn--primary btn--lg"
            disabled={loading}
            style={{ marginTop: 8 }}
          >
            {loading ? "Creando..." : "Crear evento"}
          </button>
        </form>
      </div>

      {/* Lista de eventos existentes */}
      <div className="admin-section">
        <h2 className="admin-section__title">Eventos existentes</h2>
        {fetching ? (
          <div className="loading"><div className="spinner" /> Cargando...</div>
        ) : events.length === 0 ? (
          <p className="empty__desc">No hay eventos cargados.</p>
        ) : (
          <div className="admin-table-wrap">
            <table className="admin-table">
              <thead>
                <tr>
                  <th>#</th>
                  <th>Título</th>
                  <th>Categoría</th>
                  <th>Fecha</th>
                  <th>Lugar</th>
                  <th>Precio</th>
                  <th>Cupo</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {events.map((ev) => (
                  <tr key={ev.id}>
                    <td className="admin-table__id">{ev.id}</td>
                    <td className="admin-table__title">{ev.titulo}</td>
                    <td><span className="badge badge--recitales" style={badgeStyle(ev.categoria)}>{ev.categoria}</span></td>
                    <td>{new Date(ev.fecha).toLocaleDateString("es-AR")}</td>
                    <td>{ev.lugar}</td>
                    <td>${Number(ev.precio).toLocaleString("es-AR")}</td>
                    <td>{ev.cupo_disponible}/{ev.cupo_maximo}</td>
                    <td>
                      <button
                        className="btn btn--danger btn--sm"
                        onClick={() => handleEliminar(ev)}
                      >
                        Eliminar
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}

function badgeStyle(categoria) {
  const map = {
    Recitales: { background: "#ede9fe", color: "#5b21b6" },
    Teatro:    { background: "#dcfce7", color: "#166534" },
    Deportes:  { background: "#dbeafe", color: "#1e40af" },
  };
  return map[categoria] ?? {};
}

export default AdminPanel;

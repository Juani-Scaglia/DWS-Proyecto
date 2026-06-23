import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import DatePicker, { registerLocale } from "react-datepicker";
import { es } from "date-fns/locale";
import "react-datepicker/dist/react-datepicker.css";
import { getEvents, createEvent, updateEvent, deleteEvent } from "../services/eventService";
import { getVenues, createVenue, updateVenue, deleteVenue } from "../services/venueService";

registerLocale("es", es);

const EMPTY_EVENT = { titulo: "", descripcion: "", categoria: "Recitales", precio: "", venue_id: "" };
const EMPTY_VENUE = { nombre: "", direccion: "", filas: "", columnas_por_fila: "" };

const todayStart = new Date();
todayStart.setHours(0, 0, 0, 0);

export default function AdminPanel({ user }) {
  const navigate = useNavigate();
  const [tab, setTab] = useState("events");

  const [events, setEvents] = useState([]);
  const [venues, setVenues] = useState([]);

  const [eventForm, setEventForm] = useState(EMPTY_EVENT);
  const [editingEvent, setEditingEvent] = useState(null);
  const [fechaHora, setFechaHora] = useState(null);

  const [venueForm, setVenueForm] = useState(EMPTY_VENUE);
  const [editingVenue, setEditingVenue] = useState(null);

  const [loading, setLoading] = useState(false);
  const [fetching, setFetching] = useState(true);
  const [msg, setMsg] = useState({ type: "", text: "" });

  useEffect(() => {
    if (!user || user.rol !== "admin") { navigate("/"); return; }
    reload();
  }, [user]);

  const reload = () => {
    setFetching(true);
    Promise.all([getEvents(), getVenues()])
      .then(([ev, vn]) => { setEvents(ev.data); setVenues(vn.data); })
      .catch(() => flash("error", "No se pudieron cargar los datos"))
      .finally(() => setFetching(false));
  };

  const flash = (type, text) => setMsg({ type, text });
  const clearMsg = () => setMsg({ type: "", text: "" });

  /* ── helpers ── */
  const handleEF = (e) => setEventForm({ ...eventForm, [e.target.name]: e.target.value });
  const handleVF = (e) => setVenueForm({ ...venueForm, [e.target.name]: e.target.value });
  const selectedVenue = venues.find((v) => v.id === Number(eventForm.venue_id));

  /* ══════════ EVENTOS ══════════ */

  const submitEvent = async (e) => {
    e.preventDefault();
    if (!fechaHora) { flash("error", "Seleccioná una fecha y horario."); return; }
    if (!eventForm.venue_id) { flash("error", "Seleccioná un establecimiento."); return; }
    clearMsg(); setLoading(true);
    const payload = {
      ...eventForm,
      precio: parseFloat(eventForm.precio),
      venue_id: parseInt(eventForm.venue_id),
      fecha: fechaHora.toISOString(),
    };
    try {
      if (editingEvent) {
        await updateEvent(editingEvent.id, payload);
        flash("success", `Evento "${eventForm.titulo}" actualizado.`);
      } else {
        await createEvent(payload);
        flash("success", `Evento "${eventForm.titulo}" creado.`);
      }
      setEventForm(EMPTY_EVENT); setFechaHora(null); setEditingEvent(null);
      reload();
    } catch (err) {
      flash("error", err.response?.data?.error || "Error al guardar el evento");
    } finally { setLoading(false); }
  };

  const startEditEvent = (ev) => {
    setEditingEvent(ev);
    setEventForm({
      titulo: ev.titulo, descripcion: ev.descripcion || "",
      categoria: ev.categoria, precio: String(ev.precio), venue_id: String(ev.venue_id),
    });
    setFechaHora(new Date(ev.fecha));
    setTab("events");
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const cancelEditEvent = () => { setEditingEvent(null); setEventForm(EMPTY_EVENT); setFechaHora(null); };

  const removeEvent = async (ev) => {
    if (!window.confirm(`¿Eliminar "${ev.titulo}"? Se cancelarán todas las entradas.`)) return;
    clearMsg();
    try { await deleteEvent(ev.id); flash("success", `"${ev.titulo}" eliminado.`); reload(); }
    catch (err) { flash("error", err.response?.data?.error || "Error al eliminar"); }
  };

  /* ══════════ VENUES ══════════ */

  const submitVenue = async (e) => {
    e.preventDefault(); clearMsg(); setLoading(true);
    const payload = {
      nombre: venueForm.nombre, direccion: venueForm.direccion,
      filas: parseInt(venueForm.filas), columnas_por_fila: parseInt(venueForm.columnas_por_fila),
    };
    try {
      if (editingVenue) {
        await updateVenue(editingVenue.id, payload);
        flash("success", `"${venueForm.nombre}" actualizado.`);
      } else {
        await createVenue(payload);
        flash("success", `"${venueForm.nombre}" creado.`);
      }
      setVenueForm(EMPTY_VENUE); setEditingVenue(null); reload();
    } catch (err) {
      flash("error", err.response?.data?.error || "Error al guardar el establecimiento");
    } finally { setLoading(false); }
  };

  const startEditVenue = (v) => {
    setEditingVenue(v);
    setVenueForm({ nombre: v.nombre, direccion: v.direccion, filas: String(v.filas), columnas_por_fila: String(v.columnas_por_fila) });
    setTab("venues");
    window.scrollTo({ top: 0, behavior: "smooth" });
  };

  const cancelEditVenue = () => { setEditingVenue(null); setVenueForm(EMPTY_VENUE); };

  const removeVenue = async (v) => {
    if (!window.confirm(`¿Eliminar "${v.nombre}"?`)) return;
    clearMsg();
    try { await deleteVenue(v.id); flash("success", `"${v.nombre}" eliminado.`); reload(); }
    catch (err) { flash("error", err.response?.data?.error || "Error al eliminar"); }
  };

  /* ══════════ RENDER ══════════ */

  return (
    <div className="admin-page">
      <div className="page-header">
        <h1 className="page-title">Panel de Administración</h1>
        <p className="page-subtitle">Gestioná establecimientos y eventos del sistema</p>
      </div>

      {msg.text && <div className={`alert alert--${msg.type}`}>{msg.text}</div>}

      <div className="admin-tabs">
        <button className={`admin-tab${tab === "events" ? " admin-tab--active" : ""}`} onClick={() => setTab("events")}>Eventos</button>
        <button className={`admin-tab${tab === "venues" ? " admin-tab--active" : ""}`} onClick={() => setTab("venues")}>Establecimientos</button>
      </div>

      {/* ─── TAB EVENTOS ─── */}
      {tab === "events" && (
        <>
          <div className="admin-section">
            <h2 className="admin-section__title">{editingEvent ? "Editar Evento" : "Nuevo Evento"}</h2>
            <form onSubmit={submitEvent} className="admin-form">
              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Título</label>
                  <input className="form-input" name="titulo" placeholder="Ej: Coldplay en Córdoba" value={eventForm.titulo} onChange={handleEF} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Categoría</label>
                  <select className="form-input" name="categoria" value={eventForm.categoria} onChange={handleEF} required>
                    <option value="Recitales">Recitales</option>
                    <option value="Teatro">Teatro</option>
                    <option value="Deportes">Deportes</option>
                    <option value="Cine">Cine</option>
                    <option value="Otra">Otra</option>
                  </select>
                </div>
              </div>

              <div className="form-group">
                <label className="form-label">Descripción</label>
                <textarea className="form-input form-textarea" name="descripcion" placeholder="Descripción del evento..." value={eventForm.descripcion} onChange={handleEF} rows={3} />
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Fecha</label>
                  <DatePicker locale="es" selected={fechaHora} onChange={(d) => setFechaHora(d)} minDate={todayStart} dateFormat="dd/MM/yyyy" placeholderText="Seleccioná fecha" className="form-input" calendarClassName="dp-calendar" required />
                </div>
                <div className="form-group">
                  <label className="form-label">Horario</label>
                  <DatePicker locale="es" selected={fechaHora} onChange={(d) => setFechaHora(d)} showTimeSelect showTimeSelectOnly timeIntervals={15} timeCaption="Hora" dateFormat="HH:mm" timeFormat="HH:mm" placeholderText="Horario" className="form-input" calendarClassName="dp-calendar" required />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Establecimiento *</label>
                  <select className="form-input" name="venue_id" value={eventForm.venue_id} onChange={handleEF} required>
                    <option value="">— Seleccioná —</option>
                    {venues.map((v) => (
                      <option key={v.id} value={v.id}>{v.nombre} ({v.filas}×{v.columnas_por_fila} = {v.capacidad} asientos)</option>
                    ))}
                  </select>
                  {selectedVenue && (
                    <p className="form-hint">{selectedVenue.direccion} — Capacidad: {selectedVenue.capacidad} asientos</p>
                  )}
                </div>
                <div className="form-group">
                  <label className="form-label">Precio ($)</label>
                  <input className="form-input" name="precio" type="number" min="1" step="0.01" placeholder="5000" value={eventForm.precio} onChange={handleEF} required />
                </div>
              </div>

              <div style={{ display: "flex", gap: 8, marginTop: 8 }}>
                <button type="submit" className="btn btn--primary btn--lg" disabled={loading}>
                  {loading ? "Guardando..." : editingEvent ? "Guardar cambios" : "Crear evento"}
                </button>
                {editingEvent && <button type="button" className="btn btn--ghost btn--lg" onClick={cancelEditEvent}>Cancelar</button>}
              </div>
            </form>
          </div>

          <div className="admin-section">
            <h2 className="admin-section__title">Eventos existentes</h2>
            {fetching ? <p>Cargando...</p> : events.length === 0 ? <p>No hay eventos cargados.</p> : (
              <div className="admin-table-wrap">
                <table className="admin-table">
                  <thead>
                    <tr><th>#</th><th>Título</th><th>Categoría</th><th>Fecha</th><th>Lugar</th><th>Precio</th><th>Cupo</th><th></th></tr>
                  </thead>
                  <tbody>
                    {events.map((ev) => (
                      <tr key={ev.id}>
                        <td>{ev.id}</td>
                        <td><strong>{ev.titulo}</strong></td>
                        <td>{ev.categoria}</td>
                        <td>{new Date(ev.fecha).toLocaleDateString("es-AR")}</td>
                        <td>{ev.lugar}</td>
                        <td>${Number(ev.precio).toLocaleString("es-AR")}</td>
                        <td>{ev.cupo_disponible}/{ev.cupo_maximo}</td>
                        <td style={{ display: "flex", gap: 6 }}>
                          <button className="btn btn--secondary btn--sm" onClick={() => startEditEvent(ev)}>Editar</button>
                          <button className="btn btn--danger btn--sm" onClick={() => removeEvent(ev)}>Eliminar</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </>
      )}

      {/* ─── TAB ESTABLECIMIENTOS ─── */}
      {tab === "venues" && (
        <>
          <div className="admin-section">
            <h2 className="admin-section__title">{editingVenue ? "Editar Establecimiento" : "Nuevo Establecimiento"}</h2>
            <form onSubmit={submitVenue} className="admin-form">
              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Nombre</label>
                  <input className="form-input" name="nombre" placeholder="Ej: Estadio Mario Kempes" value={venueForm.nombre} onChange={handleVF} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Dirección</label>
                  <input className="form-input" name="direccion" placeholder="Ej: Av. Cárcano s/n, Córdoba" value={venueForm.direccion} onChange={handleVF} required />
                </div>
              </div>

              <div className="form-row">
                <div className="form-group">
                  <label className="form-label">Filas</label>
                  <input className="form-input" name="filas" type="number" min="1" max="26" placeholder="10" value={venueForm.filas} onChange={handleVF} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Asientos por fila</label>
                  <input className="form-input" name="columnas_por_fila" type="number" min="1" placeholder="20" value={venueForm.columnas_por_fila} onChange={handleVF} required />
                </div>
                <div className="form-group">
                  <label className="form-label">Capacidad total</label>
                  <input className="form-input" readOnly value={venueForm.filas && venueForm.columnas_por_fila ? parseInt(venueForm.filas) * parseInt(venueForm.columnas_por_fila) : "—"} />
                </div>
              </div>

              <div style={{ display: "flex", gap: 8, marginTop: 8 }}>
                <button type="submit" className="btn btn--primary btn--lg" disabled={loading}>
                  {loading ? "Guardando..." : editingVenue ? "Guardar cambios" : "Crear establecimiento"}
                </button>
                {editingVenue && <button type="button" className="btn btn--ghost btn--lg" onClick={cancelEditVenue}>Cancelar</button>}
              </div>
            </form>
          </div>

          <div className="admin-section">
            <h2 className="admin-section__title">Establecimientos existentes</h2>
            {fetching ? <p>Cargando...</p> : venues.length === 0 ? <p>No hay establecimientos.</p> : (
              <div className="admin-table-wrap">
                <table className="admin-table">
                  <thead>
                    <tr><th>#</th><th>Nombre</th><th>Dirección</th><th>Dimensiones</th><th>Capacidad</th><th></th></tr>
                  </thead>
                  <tbody>
                    {venues.map((v) => (
                      <tr key={v.id}>
                        <td>{v.id}</td>
                        <td><strong>{v.nombre}</strong></td>
                        <td>{v.direccion}</td>
                        <td>{v.filas} filas × {v.columnas_por_fila} col</td>
                        <td>{v.capacidad}</td>
                        <td style={{ display: "flex", gap: 6 }}>
                          <button className="btn btn--secondary btn--sm" onClick={() => startEditVenue(v)}>Editar</button>
                          <button className="btn btn--danger btn--sm" onClick={() => removeVenue(v)}>Eliminar</button>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}

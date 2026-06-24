import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import DatePicker, { registerLocale } from "react-datepicker";
import { es } from "date-fns/locale";
import "react-datepicker/dist/react-datepicker.css";
import { getEvents, createEvent, updateEvent, deleteEvent } from "../services/eventService";
import { getVenues, createVenue, updateVenue, deleteVenue } from "../services/venueService";
import { getEventReport } from "../services/ticketService";
import { uploadImage } from "../services/eventService";

registerLocale("es", es);

const EMPTY_EVENT = { titulo: "", descripcion: "", categoria: "Recitales", precio: "", venue_id: "", imagen: "" };
const EMPTY_VENUE = {
  nombre: "", direccion: "", tipo: "estadio",
  capacidad: "",
  cap_platea_norte: "", cap_platea_sur: "",
  cap_tribuna_este: "", cap_tribuna_oeste: "",
  cap_campo: "",
};
const SECTORES_ESTADIO = [
  { key: "cap_platea_norte", label: "Tribuna Norte" },
  { key: "cap_platea_sur", label: "Tribuna Sur" },
  { key: "cap_tribuna_este", label: "Tribuna Este" },
  { key: "cap_tribuna_oeste", label: "Tribuna Oeste" },
  { key: "cap_campo", label: "Campo" },
];

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
  const [reports, setReports] = useState({});

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

  const loadReports = async (evts) => {
    const reps = {};
    for (const ev of evts) {
      try {
        const res = await getEventReport(ev.id);
        reps[ev.id] = res.data;
      } catch { /* skip */ }
    }
    setReports(reps);
  };

  const flash = (type, text) => setMsg({ type, text });
  const clearMsg = () => setMsg({ type: "", text: "" });

  /* ── helpers ── */
  const handleEF = (e) => setEventForm({ ...eventForm, [e.target.name]: e.target.value });

  const handleImageUpload = async (e) => {
    const file = e.target.files[0];
    if (!file) return;
    try {
      const res = await uploadImage(file);
      setEventForm((prev) => ({ ...prev, imagen: res.data.url }));
    } catch {
      flash("error", "Error al subir la imagen");
    }
  };
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
      categoria: ev.categoria, precio: String(ev.precio), venue_id: String(ev.venue_id), imagen: ev.imagen || "",
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
    const isEscenario = venueForm.tipo === "escenario";
    const payload = {
      nombre: venueForm.nombre, direccion: venueForm.direccion, tipo: venueForm.tipo,
      capacidad: isEscenario ? (parseInt(venueForm.capacidad) || 0) : 0,
      cap_tribuna_norte: isEscenario ? 0 : (parseInt(venueForm.cap_platea_norte) || 0),
      cap_tribuna_sur: isEscenario ? 0 : (parseInt(venueForm.cap_platea_sur) || 0),
      cap_tribuna_este: isEscenario ? 0 : (parseInt(venueForm.cap_tribuna_este) || 0),
      cap_tribuna_oeste: isEscenario ? 0 : (parseInt(venueForm.cap_tribuna_oeste) || 0),
      cap_campo: isEscenario ? 0 : (parseInt(venueForm.cap_campo) || 0),
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
    setVenueForm({
      nombre: v.nombre, direccion: v.direccion, tipo: v.tipo || "estadio",
      capacidad: String(v.capacidad || ""),
      cap_platea_norte: String(v.cap_platea_norte || ""),
      cap_platea_sur: String(v.cap_platea_sur || ""),
      cap_tribuna_este: String(v.cap_tribuna_este || ""),
      cap_tribuna_oeste: String(v.cap_tribuna_oeste || ""),
      cap_campo: String(v.cap_campo || ""),
    });
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
        <button className={`admin-tab${tab === "reports" ? " admin-tab--active" : ""}`} onClick={() => { setTab("reports"); loadReports(events); }}>Reportes</button>
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
                <textarea className="form-input form-textarea" name="descripcion" placeholder="Descripcion del evento..." value={eventForm.descripcion} onChange={handleEF} rows={3} />
              </div>

              <div className="form-group">
                <label className="form-label">Imagen del evento</label>
                <input className="form-input" type="file" accept="image/*" onChange={handleImageUpload} />
                {eventForm.imagen && <img src={eventForm.imagen} alt="Preview" style={{ marginTop: 8, maxHeight: 120, borderRadius: 8 }} />}
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
                    {venues
                      .filter((v) => {
                        if (eventForm.categoria === "Deportes") return (v.tipo || "estadio") === "estadio";
                        return true;
                      })
                      .map((v) => (
                        <option key={v.id} value={v.id}>{v.nombre} ({v.capacidad} asientos)</option>
                      ))}
                  </select>
                  {selectedVenue && (
                    <p className="form-hint">{selectedVenue.direccion} - {selectedVenue.capacidad} asientos</p>
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

              <div className="form-group">
                <label className="form-label">Tipo de establecimiento</label>
                <select className="form-input" name="tipo" value={venueForm.tipo} onChange={handleVF}>
                  <option value="estadio">Estadio (sectores alrededor de cancha)</option>
                  <option value="escenario">Escenario (sectores frente al escenario)</option>
                </select>
              </div>

              {venueForm.tipo === "escenario" ? (
                <>
                  <div className="form-group">
                    <label className="form-label">Capacidad total</label>
                    <input className="form-input" name="capacidad" type="number" min="1" placeholder="Ej: 500" value={venueForm.capacidad} onChange={handleVF} required />
                  </div>
                </>
              ) : (
                <>
                  <p className="form-label" style={{ marginBottom: 4 }}>Capacidad por tribuna</p>
                  <div className="form-row" style={{ gridTemplateColumns: "1fr 1fr 1fr" }}>
                    {SECTORES_ESTADIO.map(({ key, label }) => (
                      <div className="form-group" key={key}>
                        <label className="form-label">{label}</label>
                        <input className="form-input" name={key} type="number" min="0" placeholder="0" value={venueForm[key]} onChange={handleVF} />
                      </div>
                    ))}
                  </div>
                  <p className="form-hint" style={{ fontSize: 14, fontWeight: 600 }}>
                    Capacidad total: {SECTORES_ESTADIO.reduce((sum, { key }) => sum + (parseInt(venueForm[key]) || 0), 0)} asientos
                  </p>
                </>
              )}

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
                    <tr><th>#</th><th>Nombre</th><th>Tipo</th><th>Direccion</th><th>Capacidad</th><th>Sectores</th><th></th></tr>
                  </thead>
                  <tbody>
                    {venues.map((v) => (
                      <tr key={v.id}>
                        <td>{v.id}</td>
                        <td><strong>{v.nombre}</strong></td>
                        <td>{v.tipo === "escenario" ? "Escenario" : "Estadio"}</td>
                        <td>{v.direccion}</td>
                        <td>{v.capacidad}</td>
                        <td style={{ fontSize: 12 }}>
                          {v.tipo === "escenario" ? `General: ${v.capacidad}` :
                            SECTORES_ESTADIO.map(({ key, label }) => {
                              const val = v[key] || 0;
                              return val > 0 ? `${label}: ${val}` : null;
                            }).filter(Boolean).join(" | ")
                          }
                        </td>
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

      {/* ─── TAB REPORTES ─── */}
      {tab === "reports" && (
        <div className="admin-section">
          <h2 className="admin-section__title">Reportes de Ocupación</h2>
          {events.length === 0 ? <p>No hay eventos.</p> : (
            <div className="reports-grid">
              {events.map((ev) => {
                const r = reports[ev.id];
                if (!r) return <div key={ev.id} className="report-card"><p>Cargando...</p></div>;
                const pct = r.porcentaje_ocupacion?.toFixed(1) || 0;
                return (
                  <div key={ev.id} className="report-card">
                    <h3 className="report-card__title">{ev.titulo}</h3>
                    <p className="report-card__venue">{ev.lugar}</p>
                    <div className="report-card__bar-wrap">
                      <div className="report-card__bar">
                        <div className="report-card__bar-fill" style={{ width: `${pct}%` }} />
                      </div>
                      <span className="report-card__pct">{pct}%</span>
                    </div>
                    <div className="report-card__stats">
                      <div><span className="report-card__stat-val">{r.entradas_vendidas}</span><span className="report-card__stat-label">Vendidas</span></div>
                      <div><span className="report-card__stat-val">{r.entradas_canceladas}</span><span className="report-card__stat-label">Canceladas</span></div>
                      <div><span className="report-card__stat-val">{r.asientos_ocupados}/{r.asientos_totales}</span><span className="report-card__stat-label">Asientos</span></div>
                      <div><span className="report-card__stat-val">{ev.cupo_disponible}</span><span className="report-card__stat-label">Disponibles</span></div>
                    </div>
                  </div>
                );
              })}
            </div>
          )}
        </div>
      )}
    </div>
  );
}

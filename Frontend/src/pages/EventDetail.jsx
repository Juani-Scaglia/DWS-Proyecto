import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { getEventById, getEventSeats } from "../services/eventService";
import { purchaseTickets } from "../services/ticketService";
import SeatMap from "../components/SeatMap";

export default function EventDetail({ user }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [event, setEvent] = useState(null);
  const [seats, setSeats] = useState([]);
  const [selectedSeats, setSelectedSeats] = useState([]);
  const [loading, setLoading] = useState(true);
  const [buying, setBuying] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");
  const [seatMapKey, setSeatMapKey] = useState(0);

  useEffect(() => {
    Promise.all([getEventById(id), getEventSeats(id)])
      .then(([evRes, seatsRes]) => {
        setEvent(evRes.data);
        setSeats(seatsRes.data);
      })
      .catch(() => setError("Evento no encontrado"))
      .finally(() => setLoading(false));
  }, [id]);

  const handleComprar = async () => {
    if (!user) { navigate("/login"); return; }
    if (selectedSeats.length === 0) { setError("Seleccioná al menos un asiento"); return; }
    setBuying(true); setError(""); setSuccess("");
    try {
      await purchaseTickets(parseInt(id), selectedSeats);
      setSuccess(`${selectedSeats.length} entrada(s) comprada(s). Podés verlas en Mis Entradas.`);
      setSelectedSeats([]);
      setSeatMapKey((k) => k + 1);
      const [evRes, seatsRes] = await Promise.all([getEventById(id), getEventSeats(id)]);
      setEvent(evRes.data);
      setSeats(seatsRes.data);
    } catch (err) {
      setError(err.response?.data?.error || "Error al comprar");
    } finally { setBuying(false); }
  };

  if (loading) return <div className="loading">Cargando evento...</div>;
  if (!event) return <div className="container"><div className="alert alert--error">{error}</div><Link to="/">Volver</Link></div>;

  const resolveVenueType = (ev) => {
    if (ev.venue?.tipo) return ev.venue.tipo;
    if (["Teatro", "Cine"].includes(ev.categoria)) return "escenario";
    if (ev.categoria === "Deportes") return "estadio";
    return "escenario";
  };

  const fecha = new Date(event.fecha).toLocaleDateString("es-AR", { weekday: "long", day: "numeric", month: "long", year: "numeric" });
  const hora = new Date(event.fecha).toLocaleTimeString("es-AR", { hour: "2-digit", minute: "2-digit" });
  const total = selectedSeats.length * Number(event.precio);

  return (
    <div className="container" style={{ paddingTop: 32, paddingBottom: 64 }}>
      <Link to="/" className="back-link">← Volver a eventos</Link>

      {event.imagen && (
        <img src={event.imagen} alt={event.titulo} className="event-detail__img" />
      )}

      <h1 style={{ margin: "16px 0 8px" }}>{event.titulo}</h1>
      <span className="badge">{event.categoria}</span>

      <p style={{ margin: "12px 0", color: "#666" }}>{event.descripcion}</p>

      <div className="info-grid">
        <div><p className="info-item__label">Fecha</p><p>{fecha}</p></div>
        <div><p className="info-item__label">Horario</p><p>{hora} hs</p></div>
        <div><p className="info-item__label">Lugar</p><p>{event.lugar}</p></div>
        <div><p className="info-item__label">Disponibles</p><p>{event.cupo_disponible} / {event.cupo_maximo}</p></div>
      </div>

      <hr style={{ margin: "24px 0" }} />

      {success && <div className="alert alert--success">{success}</div>}
      {error && <div className="alert alert--error">{error}</div>}

      <h2>Seleccioná tus asientos</h2>
      {seats.length > 0 ? (
        <SeatMap key={seatMapKey} seats={seats} maxSelectable={10} onSelectionChange={setSelectedSeats} venueType={resolveVenueType(event)} eventCategory={event.categoria} />
      ) : (
        <p>No hay asientos disponibles.</p>
      )}

      <div className="purchase-bar">
        <div>
          <p><strong>{selectedSeats.length}</strong> asiento(s) seleccionado(s)</p>
          {selectedSeats.length > 0 && (
            <p className="purchase-bar__total">Total: ${total.toLocaleString("es-AR")}</p>
          )}
        </div>
        <button
          className="btn btn--primary btn--lg"
          onClick={handleComprar}
          disabled={buying || selectedSeats.length === 0 || event.cupo_disponible === 0}
        >
          {buying ? "Procesando..." : !user ? "Iniciá sesión para comprar" : `Comprar ${selectedSeats.length} entrada(s)`}
        </button>
      </div>
    </div>
  );
}

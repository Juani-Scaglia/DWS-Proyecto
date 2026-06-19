import { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { getEventById } from "../services/eventService";
import { purchaseTicket, getMyTickets } from "../services/ticketService";

const CAT_CLASS = {
  Recitales: "badge--recitales",
  Teatro:    "badge--teatro",
  Deportes:  "badge--deportes",
  Cine:      "badge--cine",
  Otra:      "badge--otra",
};

function EventDetail({ user }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [event, setEvent]               = useState(null);
  const [loading, setLoading]           = useState(true);
  const [buying, setBuying]             = useState(false);
  const [cantidad, setCantidad]         = useState(1);
  const [userTicketCount, setUserTicketCount] = useState(0);
  const [error, setError]               = useState("");
  const [success, setSuccess]           = useState("");

  useEffect(() => {
    getEventById(id)
      .then((res) => setEvent(res.data))
      .catch(() => setError("Evento no encontrado"))
      .finally(() => setLoading(false));
  }, [id]);

  useEffect(() => {
    if (!user) return;
    getMyTickets()
      .then((res) => {
        const count = res.data.filter(
          (t) => t.event_id === parseInt(id) && t.estado === "activo"
        ).length;
        setUserTicketCount(count);
      })
      .catch(() => {});
  }, [id, user]);

  const handleComprar = async () => {
    if (!user) { navigate("/login"); return; }
    setBuying(true);
    setError("");
    setSuccess("");
    try {
      for (let i = 0; i < cantidad; i++) {
        await purchaseTicket(parseInt(id));
      }
      const msg = cantidad > 1
        ? `${cantidad} entradas compradas. Podés verlas en Mis Entradas.`
        : "Entrada comprada. Podés verla en Mis Entradas.";
      setSuccess(msg);
      setEvent((prev) => ({ ...prev, cupo_disponible: prev.cupo_disponible - cantidad }));
      setUserTicketCount((prev) => prev + cantidad);
      setCantidad(1);
    } catch (err) {
      setError(err.response?.data?.error || "Error al comprar la entrada");
    } finally {
      setBuying(false);
    }
  };

  if (loading) {
    return (
      <div className="loading">
        <div className="spinner" /> Cargando evento...
      </div>
    );
  }

  if (error && !event) {
    return (
      <div className="event-detail-page">
        <div className="event-detail">
          <div className="alert alert--error">{error}</div>
          <Link to="/" className="btn btn--secondary">← Volver a eventos</Link>
        </div>
      </div>
    );
  }

  const fecha = new Date(event.fecha).toLocaleDateString("es-AR", {
    weekday: "long", day: "numeric", month: "long", year: "numeric",
  });

  const hora = new Date(event.fecha).toLocaleTimeString("es-AR", {
    hour: "2-digit", minute: "2-digit",
  });

  const soldOut        = event.cupo_disponible === 0;
  const atLimit        = user && userTicketCount >= 10;
  const remainingSlots = Math.max(0, 10 - userTicketCount);
  const maxCantidad    = Math.min(10, event.cupo_disponible, remainingSlots);
  const total          = Number(event.precio) * cantidad;

  let btnLabel = cantidad > 1 ? `Comprar ${cantidad} entradas` : "Comprar entrada";
  if (buying)        btnLabel = "Procesando...";
  else if (atLimit)  btnLabel = "Límite alcanzado";
  else if (soldOut)  btnLabel = "Sin cupo disponible";
  else if (!user)    btnLabel = "Iniciá sesión para comprar";

  return (
    <div className="event-detail-page" onClick={() => navigate("/")}>
      <div className="event-detail">
        <Link to="/" className="back-link" onClick={(e) => e.stopPropagation()}>← Volver a eventos</Link>

        <div className="event-detail__card" onClick={(e) => e.stopPropagation()}>
        <div className="event-detail__banner">
          <span className={`badge ${CAT_CLASS[event.categoria] ?? ""}`}>
            {event.categoria}
          </span>
          <h1 className="event-detail__title">{event.titulo}</h1>
        </div>

        <div className="event-detail__body">
          <p className="event-detail__desc">{event.descripcion}</p>

          <div className="info-grid">
            <div>
              <p className="info-item__label">Fecha</p>
              <p className="info-item__value">{fecha}</p>
            </div>
            <div>
              <p className="info-item__label">Horario</p>
              <p className="info-item__value">{hora} hs</p>
            </div>
            <div>
              <p className="info-item__label">Lugar</p>
              <p className="info-item__value">{event.lugar}</p>
            </div>
            <div>
              <p className="info-item__label">Entradas disponibles</p>
              <p className="info-item__value">
                {soldOut ? "Agotadas" : `${event.cupo_disponible} de ${event.cupo_maximo}`}
              </p>
            </div>
          </div>

          <hr className="divider" />

          {success && <div className="alert alert--success">{success}</div>}
          {error   && <div className="alert alert--error">{error}</div>}
          {atLimit && !success && (
            <div className="alert alert--warning">
              Alcanzaste el límite de 10 entradas para este evento.
            </div>
          )}

          <div className="purchase-row">
            <div>
              <p className="price-tag__label">Precio por entrada</p>
              <p className="price-tag__amount">
                ${Number(event.precio).toLocaleString("es-AR")}
              </p>

              {!soldOut && !atLimit && user && (
                <div className="qty-row">
                  <span className="qty-label">Cantidad</span>
                  <div className="qty-selector">
                    <button
                      type="button"
                      className="qty-btn"
                      onClick={() => setCantidad((c) => Math.max(1, c - 1))}
                      disabled={buying || cantidad <= 1}
                    >
                      −
                    </button>
                    <span className="qty-value">{cantidad}</span>
                    <button
                      type="button"
                      className="qty-btn"
                      onClick={() => setCantidad((c) => Math.min(maxCantidad, c + 1))}
                      disabled={buying || cantidad >= maxCantidad}
                    >
                      +
                    </button>
                  </div>
                </div>
              )}
            </div>

            <div className="purchase-action">
              {cantidad > 1 && !soldOut && !atLimit && user && (
                <p className="price-tag__total">
                  Total: ${total.toLocaleString("es-AR")}
                </p>
              )}
              <button
                className="btn btn--primary btn--lg"
                onClick={handleComprar}
                disabled={soldOut || buying || atLimit}
              >
                {btnLabel}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
  );
}

export default EventDetail;

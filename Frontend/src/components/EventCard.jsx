import { Link } from "react-router-dom";

const CAT_CLASS = {
  Recitales: "badge--recitales",
  Teatro:    "badge--teatro",
  Deportes:  "badge--deportes",
  Cine:      "badge--cine",
  Otra:      "badge--otra",
};

function EventCard({ event }) {
  const pct      = event.cupo_maximo > 0
    ? Math.round((event.cupo_disponible / event.cupo_maximo) * 100)
    : 0;
  const isSoldOut  = event.cupo_disponible === 0;
  const isCritical = !isSoldOut && event.cupo_disponible <= 10;
  const isLow      = !isSoldOut && !isCritical && pct <= 20;

  const fillClass = isCritical ? " capacity__fill--critical"
    : isLow ? " capacity__fill--low"
    : "";

  const fecha = new Date(event.fecha).toLocaleDateString("es-AR", {
    weekday: "long", day: "numeric", month: "long", year: "numeric",
  });

  return (
    <Link to={`/event/${event.id}`} className="event-card">
      <div className="event-card__header">
        <span className={`badge ${CAT_CLASS[event.categoria] ?? ""}`}>
          {event.categoria}
        </span>
        {isSoldOut && <span className="badge badge--soldout">Agotado</span>}
      </div>

      <div className="event-card__body">
        <h3 className="event-card__title">{event.titulo}</h3>
        <p className="event-card__meta">{event.lugar}</p>
        <p className="event-card__meta">{fecha}</p>

        <div className="capacity">
          <span>{isSoldOut ? "Sin cupo" : `${event.cupo_disponible} disponibles`}</span>
          <div className="capacity__bar">
            <div className={`capacity__fill${fillClass}`} style={{ width: `${pct}%` }} />
          </div>
        </div>
      </div>

      <div className="event-card__footer">
        <span className="event-card__price">
          ${Number(event.precio).toLocaleString("es-AR")}
        </span>
        <span className="btn btn--primary btn--sm">Ver detalle</span>
      </div>
    </Link>
  );
}

export default EventCard;

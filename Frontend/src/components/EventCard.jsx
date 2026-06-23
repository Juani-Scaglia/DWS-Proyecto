import { Link } from "react-router-dom";

export default function EventCard({ event }) {
  return (
    <Link to={`/event/${event.id}`} className="event-card">
      {event.imagen && (
        <img src={event.imagen} alt={event.titulo} className="event-card__img" />
      )}
      <div className="event-card__body">
        <span className="badge">{event.categoria}</span>
        <h3 className="event-card__title">{event.titulo}</h3>
        <p className="event-card__meta">{event.lugar}</p>
        <p className="event-card__meta">
          {new Date(event.fecha).toLocaleDateString("es-AR", { day: "numeric", month: "long", year: "numeric" })}
        </p>
        <div className="event-card__footer">
          <span className="event-card__price">${Number(event.precio).toLocaleString("es-AR")}</span>
          <span className="event-card__cupo">{event.cupo_disponible} disponibles</span>
        </div>
      </div>
    </Link>
  );
}

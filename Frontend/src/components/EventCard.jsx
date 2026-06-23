import { Link } from "react-router-dom";

function EventCard({ event }) {
  return (
    <div>
      <h3>{event.titulo}</h3>

      <p>{event.categoria}</p>

      <p>${event.precio}</p>

      <Link to={`/event/${event.id}`}>
        Ver Detalle
      </Link>
    </div>
  );
}

export default EventCard;
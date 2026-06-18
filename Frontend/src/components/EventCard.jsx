import { Link } from "react-router-dom";

function EventCard({ event }) {
  return (
    <div style={{ border: "1px solid #ccc", padding: "16px", margin: "8px", borderRadius: "8px" }}>
      <h3>{event.titulo}</h3>
      <p>{event.categoria} — {event.lugar}</p>
      <p>Fecha: {new Date(event.fecha).toLocaleDateString("es-AR")}</p>
      <p>Precio: ${event.precio}</p>
      <p>Cupo disponible: {event.cupo_disponible}</p>
      <Link to={`/event/${event.id}`}>Ver detalle</Link>
    </div>
  );
}

export default EventCard;

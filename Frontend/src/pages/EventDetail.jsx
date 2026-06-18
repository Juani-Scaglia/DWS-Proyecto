import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { getEventById } from "../services/eventService";
import { purchaseTicket } from "../services/ticketService";

function EventDetail({ user }) {
  const { id } = useParams();
  const navigate = useNavigate();
  const [event, setEvent] = useState(null);
  const [error, setError] = useState("");
  const [mensaje, setMensaje] = useState("");

  useEffect(() => {
    getEventById(id)
      .then((res) => setEvent(res.data))
      .catch(() => setError("Evento no encontrado"));
  }, [id]);

  const handleComprar = async () => {
    if (!user) { navigate("/login"); return; }
    try {
      await purchaseTicket(parseInt(id));
      setMensaje("¡Ticket comprado exitosamente!");
      setEvent((prev) => ({ ...prev, cupo_disponible: prev.cupo_disponible - 1 }));
    } catch (err) {
      setError(err.response?.data?.error || "Error al comprar el ticket");
    }
  };

  if (error) return <p style={{ color: "red" }}>{error}</p>;
  if (!event) return <p>Cargando...</p>;

  return (
    <div>
      <h1>{event.titulo}</h1>
      <p>{event.descripcion}</p>
      <p><strong>Categoría:</strong> {event.categoria}</p>
      <p><strong>Fecha:</strong> {new Date(event.fecha).toLocaleDateString("es-AR")}</p>
      <p><strong>Lugar:</strong> {event.lugar}</p>
      <p><strong>Precio:</strong> ${event.precio}</p>
      <p><strong>Cupo disponible:</strong> {event.cupo_disponible}</p>
      {mensaje && <p style={{ color: "green" }}>{mensaje}</p>}
      <button onClick={handleComprar} disabled={event.cupo_disponible === 0}>
        {event.cupo_disponible === 0 ? "Sin cupo" : "Comprar ticket"}
      </button>
    </div>
  );
}

export default EventDetail;

import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";

import { getEventById } from "../services/eventService";

function EventDetail() {
  const { id } = useParams();

  //const [event, setEvent] = useState(null);

  const [event, setEvent] = useState({
    id: 1,
    titulo: "Cosquín Rock",
    descripcion: "Festival de rock",
    categoria: "Recitales",
    fecha: "2025-12-01",
    lugar: "Córdoba",
    precio: 5000,
    cupo_disponible: 80,
  });

  const loadEvent = async () => {
    try {
      const data = await getEventById(id);

      setEvent(data);
    } catch (error) {
      console.error("Error al obtener el evento:", error);
    }
  };

  /*useEffect(() => {
    loadEvent();
  }, []);*/

  if (!event) {
    return <h2>Cargando evento...</h2>;
  }

  return (
    <div>
      <h1>{event.titulo}</h1>

      <p>
        <strong>Descripción:</strong> {event.descripcion}
      </p>

      <p>
        <strong>Categoría:</strong> {event.categoria}
      </p>

      <p>
        <strong>Fecha:</strong> {event.fecha}
      </p>

      <p>
        <strong>Lugar:</strong> {event.lugar}
      </p>

      <p>
        <strong>Precio:</strong> ${event.precio}
      </p>

      <p>
        <strong>Cupo disponible:</strong> {event.cupo_disponible}
      </p>
      <button>Comprar Entrada</button>
    </div>
  );
}

export default EventDetail;

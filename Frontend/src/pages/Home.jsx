import { useEffect, useState } from "react";

import EventCard from "../components/EventCard";
import { getEvents } from "../services/eventService";

function Home() {
  const [events, setEvents] = useState([]);

  const loadEvents = async () => {
    try {
      const data = await getEvents();
      setEvents(data);
    } catch (error) {
      console.error("Error al obtener eventos:", error);
    }
  };

  /*useEffect(() => {
    loadEvents();
  }, []);*/

  useEffect(() => {
    setEvents([
      {
        id: 1,
        titulo: "Cosquín Rock",
        categoria: "Recitales",
        precio: 5000,
      },
      {
        id: 2,
        titulo: "Coldplay",
        categoria: "Recitales",
        precio: 10000,
      },
    ]);
  }, []);

  return (
    <div>
      <h1>Catálogo de Eventos</h1>

      {events.length === 0 ? (
        <p>No hay eventos disponibles.</p>
      ) : (
        events.map((event) => <EventCard key={event.id} event={event} />)
      )}
    </div>
  );
}

export default Home;

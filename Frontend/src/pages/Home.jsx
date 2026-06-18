import { useEffect, useState } from "react";
import { getEvents } from "../services/eventService";
import EventCard from "../components/EventCard";

const CATEGORIAS = ["", "Recitales", "Teatro", "Deportes"];

function Home() {
  const [events, setEvents] = useState([]);
  const [category, setCategory] = useState("");
  const [error, setError] = useState("");

  useEffect(() => {
    getEvents(category)
      .then((res) => setEvents(res.data))
      .catch(() => setError("No se pudieron cargar los eventos"));
  }, [category]);

  return (
    <div>
      <h1>Catálogo de Eventos</h1>
      <div>
        {CATEGORIAS.map((cat) => (
          <button
            key={cat}
            onClick={() => setCategory(cat)}
            style={{ fontWeight: category === cat ? "bold" : "normal", margin: "4px" }}
          >
            {cat || "Todos"}
          </button>
        ))}
      </div>
      {error && <p style={{ color: "red" }}>{error}</p>}
      <div>
        {events.map((event) => (
          <EventCard key={event.id} event={event} />
        ))}
      </div>
    </div>
  );
}

export default Home;

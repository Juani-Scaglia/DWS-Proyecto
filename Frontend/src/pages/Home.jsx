import { useEffect, useState } from "react";
import { getEvents } from "../services/eventService";
import EventCard from "../components/EventCard";

const CATEGORIAS = [
  { value: "", label: "Todos" },
  { value: "Recitales", label: "Recitales" },
  { value: "Teatro", label: "Teatro" },
  { value: "Deportes", label: "Deportes" },
  { value: "Cine", label: "Cine" },
  { value: "Otra", label: "Otra" },
];

export default function Home() {
  const [events, setEvents] = useState([]);
  const [category, setCategory] = useState("");
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    setLoading(true);
    getEvents(category)
      .then((res) => { setEvents(res.data); setError(""); })
      .catch(() => setError("No se pudieron cargar los eventos."))
      .finally(() => setLoading(false));
  }, [category]);

  return (
    <div className="container" style={{ paddingTop: 32, paddingBottom: 64 }}>
      <h1>Catálogo de Eventos</h1>

      <div className="filters">
        {CATEGORIAS.map(({ value, label }) => (
          <button
            key={value}
            className={`filter-btn${category === value ? " filter-btn--active" : ""}`}
            onClick={() => setCategory(value)}
          >
            {label}
          </button>
        ))}
      </div>

      {error && <div className="alert alert--error">{error}</div>}

      {loading ? (
        <p>Cargando eventos...</p>
      ) : events.length === 0 ? (
        <p>No hay eventos en esta categoría.</p>
      ) : (
        <div className="events-grid">
          {events.map((event) => (
            <EventCard key={event.id} event={event} />
          ))}
        </div>
      )}
    </div>
  );
}

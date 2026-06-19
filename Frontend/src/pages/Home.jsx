import { useEffect, useState } from "react";
import { getEvents } from "../services/eventService";
import EventCard from "../components/EventCard";

const CATEGORIAS = [
  { value: "",          label: "Todos" },
  { value: "Recitales", label: "Recitales" },
  { value: "Teatro",    label: "Teatro" },
  { value: "Deportes",  label: "Deportes" },
  { value: "Cine",      label: "Cine" },
  { value: "Otra",      label: "Otra" },
];

function Home() {
  const [events, setEvents]     = useState([]);
  const [category, setCategory] = useState("");
  const [loading, setLoading]   = useState(true);
  const [error, setError]       = useState("");

  useEffect(() => {
    setLoading(true);
    getEvents(category)
      .then((res) => { setEvents(res.data); setError(""); })
      .catch(() => setError("No se pudieron cargar los eventos. Verificá que el servidor esté activo."))
      .finally(() => setLoading(false));
  }, [category]);

  return (
    <>
      <div className="hero">
        <h1 className="hero__title">Encontrá tu próximo evento</h1>
        <p className="hero__subtitle">Recitales, teatro, deportes y más — todo en un lugar</p>
      </div>

      <div className="container">
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
          <div className="loading">
            <div className="spinner" />
            Cargando eventos...
          </div>
        ) : events.length === 0 && !error ? (
          <div className="empty">
            <p className="empty__title">No hay eventos en esta categoría</p>
            <p className="empty__desc">Probá con otra categoría o volvé más tarde.</p>
          </div>
        ) : (
          <div className="events-grid">
            {events.map((event) => (
              <EventCard key={event.id} event={event} />
            ))}
          </div>
        )}
      </div>
    </>
  );
}

export default Home;

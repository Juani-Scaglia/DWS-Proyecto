import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { getMyTickets } from "../services/ticketService";
import TicketCard from "../components/TicketCard";

function MyTickets({ user }) {
  const navigate = useNavigate();
  const [tickets, setTickets]   = useState([]);
  const [error, setError]       = useState("");
  const [loading, setLoading]   = useState(true);
  const [pageMsg, setPageMsg]   = useState({ type: "", text: "" });

  useEffect(() => {
    if (!user) { navigate("/login"); return; }
    setLoading(true);
    getMyTickets()
      .then((res) => { setTickets(res.data); setError(""); })
      .catch(() => setError("No se pudieron cargar las entradas"))
      .finally(() => setLoading(false));
  }, [user]);

  const handleRemove = (id, message) => {
    setTickets((prev) => prev.filter((t) => t.id !== id));
    setPageMsg({ type: "success", text: message });
  };

  const handleError = (message) => {
    setPageMsg({ type: "error", text: message });
  };

  return (
    <div className="tickets-page">
      <div className="page-header">
        <h1 className="page-title">Mis Entradas</h1>
        <p className="page-subtitle">Gestioná tus entradas: cancelá o transferilas a otro usuario</p>
      </div>

      {pageMsg.text && (
        <div className={`alert alert--${pageMsg.type}`}>
          {pageMsg.text}
        </div>
      )}

      {error && <div className="alert alert--error">{error}</div>}

      {loading ? (
        <div className="loading">
          <div className="spinner" /> Cargando entradas...
        </div>
      ) : tickets.length === 0 && !error ? (
        <div className="empty">
          <p className="empty__title">No tenés entradas todavía</p>
          <p className="empty__desc">Explorá los eventos disponibles y comprá tu primera entrada.</p>
          <Link to="/" className="btn btn--primary">Ver eventos</Link>
        </div>
      ) : (
        <div className="tickets-list">
          {tickets.map((ticket) => (
            <TicketCard
              key={ticket.id}
              ticket={ticket}
              onRemove={handleRemove}
              onError={handleError}
            />
          ))}
        </div>
      )}
    </div>
  );
}

export default MyTickets;

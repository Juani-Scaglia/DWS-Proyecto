import { useEffect, useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { getMyTickets } from "../services/ticketService";
import TicketCard from "../components/TicketCard";

export default function MyTickets({ user }) {
  const navigate = useNavigate();
  const [tickets, setTickets] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [msg, setMsg] = useState({ type: "", text: "" });

  useEffect(() => {
    if (!user) { navigate("/login"); return; }
    setLoading(true);
    getMyTickets()
      .then((res) => setTickets(res.data.filter((t) => t.estado === "activo")))
      .catch(() => setError("No se pudieron cargar las entradas"))
      .finally(() => setLoading(false));
  }, [user]);

  const handleRemove = (id, message) => {
    setTickets((prev) => prev.filter((t) => t.id !== id));
    setMsg({ type: "success", text: message });
  };

  return (
    <div className="container" style={{ paddingTop: 32, paddingBottom: 64 }}>
      <h1>Mis Entradas</h1>
      <p style={{ color: "#666" }}>Cancelá o transferí tus entradas a otro usuario.</p>

      {msg.text && <div className={`alert alert--${msg.type}`}>{msg.text}</div>}
      {error && <div className="alert alert--error">{error}</div>}

      {loading ? (
        <p>Cargando entradas...</p>
      ) : tickets.length === 0 ? (
        <div>
          <p>No tenés entradas todavía.</p>
          <Link to="/" className="btn btn--primary">Ver eventos</Link>
        </div>
      ) : (
        <div className="tickets-list">
          {tickets.map((ticket) => (
            <TicketCard key={ticket.id} ticket={ticket} onRemove={handleRemove} onError={(m) => setMsg({ type: "error", text: m })} />
          ))}
        </div>
      )}
    </div>
  );
}

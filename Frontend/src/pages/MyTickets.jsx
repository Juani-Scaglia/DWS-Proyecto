import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getMyTickets } from "../services/ticketService";
import TicketCard from "../components/TicketCard";

function MyTickets({ user }) {
  const navigate = useNavigate();
  const [tickets, setTickets] = useState([]);
  const [error, setError] = useState("");

  const cargarTickets = () => {
    getMyTickets()
      .then((res) => setTickets(res.data))
      .catch(() => setError("No se pudieron cargar los tickets"));
  };

  useEffect(() => {
    if (!user) { navigate("/login"); return; }
    cargarTickets();
  }, [user]);

  return (
    <div>
      <h1>Mis Entradas</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {tickets.length === 0 && !error && <p>No tenés tickets todavía.</p>}
      {tickets.map((ticket) => (
        <TicketCard key={ticket.id} ticket={ticket} onUpdate={cargarTickets} />
      ))}
    </div>
  );
}

export default MyTickets;

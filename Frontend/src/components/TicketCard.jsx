import { useState } from "react";
import { cancelTicket, transferTicket } from "../services/ticketService";

function TicketCard({ ticket, onUpdate }) {
  const [dniDestino, setDniDestino] = useState("");
  const [mensaje, setMensaje] = useState("");
  const [error, setError] = useState("");

  const handleCancelar = async () => {
    try {
      await cancelTicket(ticket.id);
      setMensaje("Ticket cancelado");
      onUpdate();
    } catch (err) {
      setError(err.response?.data?.error || "Error al cancelar");
    }
  };

  const handleTransferir = async () => {
    if (!dniDestino) return;
    try {
      await transferTicket(ticket.id, dniDestino);
      setMensaje("Ticket transferido");
      setDniDestino("");
      onUpdate();
    } catch (err) {
      setError(err.response?.data?.error || "Error al transferir");
    }
  };

  return (
    <div style={{ border: "1px solid #ccc", padding: "16px", margin: "8px", borderRadius: "8px" }}>
      <h3>Ticket #{ticket.id}</h3>
      <p>Evento: {ticket.event?.titulo || `ID ${ticket.event_id}`}</p>
      <p>Estado: {ticket.estado}</p>
      {mensaje && <p style={{ color: "green" }}>{mensaje}</p>}
      {error && <p style={{ color: "red" }}>{error}</p>}
      {ticket.estado === "activo" && (
        <>
          <button onClick={handleCancelar}>Cancelar</button>
          <br />
          <input
            placeholder="DNI destinatario"
            value={dniDestino}
            onChange={(e) => setDniDestino(e.target.value)}
          />
          <button onClick={handleTransferir}>Transferir</button>
        </>
      )}
    </div>
  );
}

export default TicketCard;

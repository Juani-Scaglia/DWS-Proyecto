import { useState } from "react";
import { cancelTicket, transferTicket } from "../services/ticketService";

export default function TicketCard({ ticket, onRemove, onError }) {
  const [dni, setDni] = useState("");
  const [loading, setLoading] = useState(false);

  const handleCancelar = async () => {
    if (!window.confirm("¿Cancelar esta entrada? No se puede deshacer.")) return;
    setLoading(true);
    try {
      await cancelTicket(ticket.id);
      onRemove(ticket.id, "Entrada cancelada correctamente.");
    } catch (err) {
      onError(err.response?.data?.error || "Error al cancelar");
    } finally { setLoading(false); }
  };

  const handleTransferir = async () => {
    if (!dni.trim()) return;
    setLoading(true);
    try {
      await transferTicket(ticket.id, dni.trim());
      onRemove(ticket.id, `Entrada transferida al DNI ${dni.trim()}.`);
    } catch (err) {
      onError(err.response?.data?.error || "Error al transferir");
    } finally { setLoading(false); }
  };

  const seatLabel = ticket.seat ? `${ticket.seat.fila}${ticket.seat.numero}` : "—";

  return (
    <div className="ticket-card">
      <div className="ticket-card__top">
        <div>
          <p className="ticket-card__event">{ticket.event?.titulo || `Evento #${ticket.event_id}`}</p>
          <p className="ticket-card__meta">Asiento: <strong>{seatLabel}</strong></p>
        </div>
        <span className="badge badge--active">Activo</span>
      </div>

      <div className="ticket-card__actions">
        <button className="btn btn--danger btn--sm" onClick={handleCancelar} disabled={loading}>
          {loading ? "Procesando..." : "Cancelar entrada"}
        </button>
        <div className="transfer-row">
          <input className="form-input" placeholder="DNI del destinatario" value={dni} onChange={(e) => setDni(e.target.value)} disabled={loading} />
          <button className="btn btn--secondary btn--sm" onClick={handleTransferir} disabled={loading || !dni.trim()}>
            Transferir
          </button>
        </div>
      </div>
    </div>
  );
}

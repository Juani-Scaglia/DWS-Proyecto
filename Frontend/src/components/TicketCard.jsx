import { useState } from "react";
import { cancelTicket, transferTicket } from "../services/ticketService";

const ESTADO_BADGE = {
  activo:      "badge--active",
  cancelado:   "badge--cancelled",
  transferido: "badge--transferred",
};

const ESTADO_LABEL = {
  activo:      "Activo",
  cancelado:   "Cancelado",
  transferido: "Transferido",
};

function TicketCard({ ticket, onRemove, onError }) {
  const [dniDestino, setDniDestino] = useState("");
  const [error, setError]           = useState("");
  const [loading, setLoading]       = useState(false);

  const handleCancelar = async () => {
    if (!window.confirm("¿Confirmás que querés cancelar esta entrada? Esta acción no se puede deshacer.")) return;
    setLoading(true);
    setError("");
    try {
      await cancelTicket(ticket.id);
      onRemove(ticket.id, "Entrada cancelada correctamente.");
    } catch (err) {
      const msg = err.response?.data?.error || "Error al cancelar la entrada";
      setError(msg);
      onError(msg);
    } finally {
      setLoading(false);
    }
  };

  const handleTransferir = async () => {
    if (!dniDestino.trim()) return;
    setLoading(true);
    setError("");
    try {
      await transferTicket(ticket.id, dniDestino.trim());
      onRemove(ticket.id, `Entrada transferida correctamente al DNI ${dniDestino.trim()}.`);
    } catch (err) {
      const msg = err.response?.data?.error || "No se encontró ningún usuario con ese DNI";
      setError(msg);
      onError(msg);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="ticket-card">
      <div className="ticket-card__top">
        <div>
          <p className="ticket-card__num">Entrada #{ticket.id}</p>
          <p className="ticket-card__event">
            {ticket.event?.titulo || `Evento ID ${ticket.event_id}`}
          </p>
        </div>
        <span className={`badge ${ESTADO_BADGE[ticket.estado] ?? ""}`}>
          {ESTADO_LABEL[ticket.estado] ?? ticket.estado}
        </span>
      </div>

      {error && <div className="alert alert--error" style={{ margin: "0 24px 12px" }}>{error}</div>}

      {ticket.estado === "activo" && (
        <div className="ticket-card__actions">
          <button
            className="btn btn--danger btn--sm"
            onClick={handleCancelar}
            disabled={loading}
          >
            {loading ? "Procesando..." : "Cancelar entrada"}
          </button>

          <div className="transfer-row">
            <input
              className="form-input"
              placeholder="DNI del destinatario"
              value={dniDestino}
              onChange={(e) => setDniDestino(e.target.value)}
              disabled={loading}
            />
            <button
              className="btn btn--secondary btn--sm"
              onClick={handleTransferir}
              disabled={loading || !dniDestino.trim()}
            >
              {loading ? "Procesando..." : "Transferir"}
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default TicketCard;

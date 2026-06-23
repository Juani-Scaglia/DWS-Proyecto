import { useState } from "react";

export default function SeatMap({ seats, maxSelectable, onSelectionChange }) {
  const [selected, setSelected] = useState([]);

  const rows = {};
  seats.forEach((s) => {
    if (!rows[s.fila]) rows[s.fila] = [];
    rows[s.fila].push(s);
  });

  const toggle = (seat) => {
    if (seat.ocupado) return;
    setSelected((prev) => {
      let next;
      if (prev.includes(seat.id)) {
        next = prev.filter((id) => id !== seat.id);
      } else {
        if (prev.length >= maxSelectable) return prev;
        next = [...prev, seat.id];
      }
      onSelectionChange(next);
      return next;
    });
  };

  return (
    <div className="seatmap">
      <div className="seatmap__stage">ESCENARIO</div>
      <div className="seatmap__grid">
        {Object.entries(rows).map(([fila, asientos]) => (
          <div key={fila} className="seatmap__row">
            <span className="seatmap__row-label">{fila}</span>
            {asientos.map((s) => {
              let cls = "seatmap__seat";
              if (s.ocupado) cls += " seatmap__seat--taken";
              else if (selected.includes(s.id)) cls += " seatmap__seat--selected";
              else cls += " seatmap__seat--free";

              return (
                <button
                  key={s.id}
                  className={cls}
                  onClick={() => toggle(s)}
                  disabled={s.ocupado}
                  title={`${s.fila}${s.numero}${s.ocupado ? " (ocupado)" : ""}`}
                >
                  {s.numero}
                </button>
              );
            })}
          </div>
        ))}
      </div>
      <div className="seatmap__legend">
        <span><span className="seatmap__dot seatmap__dot--free" /> Libre</span>
        <span><span className="seatmap__dot seatmap__dot--selected" /> Seleccionado</span>
        <span><span className="seatmap__dot seatmap__dot--taken" /> Ocupado</span>
      </div>
    </div>
  );
}

import { useState, useMemo } from "react";

export default function SeatMap({ seats, maxSelectable, onSelectionChange, venueType, eventCategory }) {
  const [selected, setSelected] = useState([]);
  const [activeSector, setActiveSector] = useState(null);

  const sectors = useMemo(() => {
    const map = {};
    seats.forEach((s) => {
      if (!map[s.sector]) map[s.sector] = [];
      map[s.sector].push(s);
    });
    return Object.entries(map).map(([nombre, asientos]) => ({
      nombre,
      seats: asientos,
      total: asientos.length,
      libres: asientos.filter((s) => !s.ocupado).length,
    }));
  }, [seats]);

  const hasSectors = sectors.length > 1;
  const isEscenario = venueType === "escenario";

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

  const renderSeats = (seatList) => {
    const rows = {};
    seatList.forEach((s) => {
      if (!rows[s.fila]) rows[s.fila] = [];
      rows[s.fila].push(s);
    });

    const maxCols = Math.max(...Object.values(rows).map((r) => r.length), 1);
    const total = seatList.length;
    let sizeClass = "";
    if (total > 200 || maxCols > 20) sizeClass = " seatmap__grid--xs";
    else if (total > 80 || maxCols > 12) sizeClass = " seatmap__grid--sm";

    return (
      <div className={`seatmap__grid${sizeClass}`}>
        {Object.entries(rows)
          .sort(([a], [b]) => a.localeCompare(b))
          .map(([fila, asientos]) => (
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
                    title={`${s.sector} - ${s.fila}${s.numero}`}
                  >
                    {sizeClass === " seatmap__grid--xs" ? "" : s.numero}
                  </button>
                );
              })}
            </div>
          ))}
      </div>
    );
  };

  // ── Escenario: siempre grilla directa, sin sectores ──
  if (isEscenario) {
    return (
      <div className="seatmap">
        <div className="seatmap__stage">ESCENARIO</div>
        {renderSeats(seats)}
        <Legend />
      </div>
    );
  }

  // ── Un solo sector: mapa directo ──
  if (!hasSectors) {
    return (
      <div className="seatmap">
        <div className="seatmap__stage">{sectors[0]?.nombre || "GENERAL"}</div>
        {renderSeats(seats)}
        <Legend />
      </div>
    );
  }

  // ── Vista detalle de un sector seleccionado ──
  if (activeSector) {
    const sec = sectors.find((s) => s.nombre === activeSector);
    return (
      <div className="seatmap">
        <button
          className="btn btn--secondary btn--sm"
          style={{ marginBottom: 12 }}
          onClick={() => setActiveSector(null)}
        >
          ← Volver al mapa
        </button>
        <div className="seatmap__stage">{activeSector}</div>
        {renderSeats(sec.seats)}
        <Legend />
      </div>
    );
  }

  const selInSec = (sec) =>
    sec ? selected.filter((id) => sec.seats.some((s) => s.id === id)).length : 0;

  // Busca sectores por patron (funciona con "Platea Norte", "Tribuna Norte", etc.)
  const findByPattern = (pattern) => sectors.find((s) => s.nombre.toLowerCase().includes(pattern));

  // ── ESTADIO: sectores alrededor de la cancha ──
  const norte = findByPattern("norte");
  const sur = findByPattern("sur");
  const este = findByPattern("este");
  const oeste = findByPattern("oeste");
  const campo = findByPattern("campo");
  const hideCampo = eventCategory === "Deportes";

  return (
    <div className="seatmap">
      <p style={{ textAlign: "center", marginBottom: 12, fontSize: 13, color: "var(--text-muted)" }}>
        Selecciona un sector para elegir tus asientos
      </p>
      <div className="stadium">
        <SectorBtn sec={norte} pos="top" selCount={selInSec(norte)} onClick={setActiveSector} />
        <div className="stadium__middle">
          <SectorBtn sec={oeste} pos="left" selCount={selInSec(oeste)} onClick={setActiveSector} />
          <div className="stadium__field">
            {!hideCampo && campo
              ? <SectorBtn sec={campo} pos="field" selCount={selInSec(campo)} onClick={setActiveSector} />
              : <span className="stadium__field-label">CANCHA</span>
            }
          </div>
          <SectorBtn sec={este} pos="right" selCount={selInSec(este)} onClick={setActiveSector} />
        </div>
        <SectorBtn sec={sur} pos="bottom" selCount={selInSec(sur)} onClick={setActiveSector} />
      </div>
    </div>
  );
}

function SectorBtn({ sec, pos, selCount, onClick }) {
  if (!sec) return null;
  const pctOcupado = sec.total > 0 ? Math.round(((sec.total - sec.libres) / sec.total) * 100) : 0;
  return (
    <button className={`sector-btn sector-btn--${pos}`} onClick={() => onClick(sec.nombre)}>
      <span className="sector-btn__name">{sec.nombre}</span>
      <span className="sector-btn__info">{sec.libres} / {sec.total} libres</span>
      <div className="sector-btn__bar">
        <div className="sector-btn__bar-fill" style={{ width: `${pctOcupado}%` }} />
      </div>
      {selCount > 0 && <span className="sector-btn__badge">{selCount} sel.</span>}
    </button>
  );
}

function Legend() {
  return (
    <div className="seatmap__legend">
      <span><span className="seatmap__dot seatmap__dot--free" /> Libre</span>
      <span><span className="seatmap__dot seatmap__dot--selected" /> Seleccionado</span>
      <span><span className="seatmap__dot seatmap__dot--taken" /> Ocupado</span>
    </div>
  );
}

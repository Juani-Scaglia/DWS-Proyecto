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

    return (
      <div className="seatmap__grid">
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
                    {s.numero}
                  </button>
                );
              })}
            </div>
          ))}
      </div>
    );
  };

  // ── Sin sectores: mapa directo ──
  if (!hasSectors) {
    return (
      <div className="seatmap">
        <div className="seatmap__stage">ESCENARIO</div>
        {renderSeats(seats)}
        <Legend />
      </div>
    );
  }

  // ── Con sectores: vista detalle de un sector ──
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

  const findSec = (name) => sectors.find((s) => s.nombre === name);
  const selInSec = (sec) =>
    sec ? selected.filter((id) => sec.seats.some((s) => s.id === id)).length : 0;

  // ── ESCENARIO: sectores apilados frente al escenario ──
  if (isEscenario) {
    const sectorOrder = ["Preferencial", "Platea Norte", "Platea Sur", "Tribuna Este", "Tribuna Oeste", "Campo"];
    const orderedSectors = sectorOrder.map(findSec).filter(Boolean);

    return (
      <div className="seatmap">
        <p style={{ textAlign: "center", marginBottom: 12, fontSize: 13, color: "var(--text-muted)" }}>
          Selecciona un sector para elegir tus asientos
        </p>
        <div className="stage-layout">
          <div className="stage-layout__stage">ESCENARIO</div>
          <div className="stage-layout__sectors">
            {orderedSectors.map((sec) => (
              <SectorBtn
                key={sec.nombre}
                sec={sec}
                pos="stage-row"
                selCount={selInSec(sec)}
                onClick={setActiveSector}
              />
            ))}
          </div>
        </div>
      </div>
    );
  }

  // ── ESTADIO: sectores alrededor de la cancha ──
  const hideCampo = eventCategory === "Deportes";
  return (
    <div className="seatmap">
      <p style={{ textAlign: "center", marginBottom: 12, fontSize: 13, color: "var(--text-muted)" }}>
        Selecciona un sector para elegir tus asientos
      </p>
      <div className="stadium">
        <SectorBtn sec={findSec("Tribuna Norte")} pos="top" selCount={selInSec(findSec("Tribuna Norte"))} onClick={setActiveSector} />
        <div className="stadium__middle">
          <SectorBtn sec={findSec("Tribuna Oeste")} pos="left" selCount={selInSec(findSec("Tribuna Oeste"))} onClick={setActiveSector} />
          <div className="stadium__field">
            {!hideCampo && <SectorBtn sec={findSec("Campo")} pos="field" selCount={selInSec(findSec("Campo"))} onClick={setActiveSector} />}
          </div>
          <SectorBtn sec={findSec("Tribuna Este")} pos="right" selCount={selInSec(findSec("Tribuna Este"))} onClick={setActiveSector} />
        </div>
        <SectorBtn sec={findSec("Tribuna Sur")} pos="bottom" selCount={selInSec(findSec("Tribuna Sur"))} onClick={setActiveSector} />
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

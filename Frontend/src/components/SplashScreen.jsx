import { useEffect, useRef, useState } from "react";

const PARTICLES = [
  { left: "8%",  bottom: "5%",  width: "5px", height: "5px",  "--delay": "0s",   "--dur": "7s"   },
  { left: "15%", bottom: "15%", width: "3px", height: "3px",  "--delay": "1.2s", "--dur": "5.5s" },
  { left: "22%", bottom: "8%",  width: "7px", height: "7px",  "--delay": "0.4s", "--dur": "6.5s" },
  { left: "30%", bottom: "20%", width: "4px", height: "4px",  "--delay": "2.1s", "--dur": "5s"   },
  { left: "38%", bottom: "3%",  width: "6px", height: "6px",  "--delay": "0.8s", "--dur": "7.5s" },
  { left: "45%", bottom: "12%", width: "3px", height: "3px",  "--delay": "1.6s", "--dur": "6s"   },
  { left: "52%", bottom: "25%", width: "5px", height: "5px",  "--delay": "0.2s", "--dur": "5.5s" },
  { left: "60%", bottom: "7%",  width: "4px", height: "4px",  "--delay": "2.8s", "--dur": "7s"   },
  { left: "68%", bottom: "18%", width: "6px", height: "6px",  "--delay": "1s",   "--dur": "6.5s" },
  { left: "75%", bottom: "10%", width: "3px", height: "3px",  "--delay": "0.6s", "--dur": "5s"   },
  { left: "82%", bottom: "4%",  width: "7px", height: "7px",  "--delay": "1.9s", "--dur": "7.5s" },
  { left: "88%", bottom: "22%", width: "4px", height: "4px",  "--delay": "0.3s", "--dur": "6s"   },
  { left: "92%", bottom: "14%", width: "5px", height: "5px",  "--delay": "2.4s", "--dur": "5.5s" },
  { left: "5%",  bottom: "35%", width: "3px", height: "3px",  "--delay": "1.4s", "--dur": "7s"   },
  { left: "95%", bottom: "28%", width: "4px", height: "4px",  "--delay": "0.9s", "--dur": "6s"   },
  { left: "50%", bottom: "2%",  width: "6px", height: "6px",  "--delay": "2s",   "--dur": "8s"   },
  { left: "3%",  bottom: "50%", width: "4px", height: "4px",  "--delay": "3.2s", "--dur": "6s"   },
  { left: "97%", bottom: "45%", width: "5px", height: "5px",  "--delay": "1.7s", "--dur": "7s"   },
];

function SplashScreen({ onDismiss }) {
  const [leaving, setLeaving] = useState(false);
  const dismissed = useRef(false);

  useEffect(() => {
    const dismiss = () => {
      if (dismissed.current) return;
      dismissed.current = true;
      setLeaving(true);
      setTimeout(onDismiss, 700);
    };
    window.addEventListener("keydown", dismiss);
    window.addEventListener("click", dismiss);
    window.addEventListener("touchstart", dismiss);
    return () => {
      window.removeEventListener("keydown", dismiss);
      window.removeEventListener("click", dismiss);
      window.removeEventListener("touchstart", dismiss);
    };
  }, [onDismiss]);

  return (
    <div className={`splash${leaving ? " splash--leaving" : ""}`}>
      <div className="splash__particles">
        {PARTICLES.map((p, i) => (
          <span key={i} className="splash__particle" style={p} />
        ))}
      </div>

      <div className="splash__rings">
        <div className="splash__ring" style={{ width: "320px", height: "320px", "--delay": "0s",   "--dur": "4s" }} />
        <div className="splash__ring" style={{ width: "540px", height: "540px", "--delay": "0.8s", "--dur": "5s" }} />
        <div className="splash__ring" style={{ width: "760px", height: "760px", "--delay": "1.6s", "--dur": "6s" }} />
      </div>

      <div className="splash__content">
        <div className="splash__monogram">TH</div>
        <h1 className="splash__title">TicketHub</h1>
        <p className="splash__tagline">Tu plataforma de eventos</p>
        <p className="splash__hint">Presioná cualquier tecla para continuar</p>
      </div>
    </div>
  );
}

export default SplashScreen;

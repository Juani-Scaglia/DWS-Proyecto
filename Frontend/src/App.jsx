import { useState } from "react";
import Navbar from "./components/Navbar";
import AppRouter from "./routes/AppRouter";
import SplashScreen from "./components/SplashScreen";

function App() {
  const [splashDone, setSplashDone] = useState(
    () => sessionStorage.getItem("splashSeen") === "1"
  );

  const [user, setUser] = useState(() => {
    const stored = localStorage.getItem("user");
    return stored ? JSON.parse(stored) : null;
  });

  const handleSplashDone = () => {
    sessionStorage.setItem("splashSeen", "1");
    setSplashDone(true);
  };

  const handleLogin = (userData, token) => {
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(userData));
    setUser(userData);
  };

  const handleLogout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  };

  return (
    <>
      {!splashDone && <SplashScreen onDismiss={handleSplashDone} />}
      <Navbar user={user} onLogout={handleLogout} />
      <AppRouter user={user} onLogin={handleLogin} />
    </>
  );
}

export default App;

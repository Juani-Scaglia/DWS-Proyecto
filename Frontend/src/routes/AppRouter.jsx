import { Routes, Route } from "react-router-dom";

import Home from "../pages/Home";
import EventDetail from "../pages/EventDetail";
import MyTickets from "../pages/MyTickets";
import Login from "../pages/Login";
import Register from "../pages/Register";
import AdminPanel from "../pages/AdminPanel";

function AppRouter({ user, onLogin }) {
  return (
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/event/:id" element={<EventDetail user={user} />} />
      <Route path="/tickets" element={<MyTickets user={user} />} />
      <Route path="/login" element={<Login onLogin={onLogin} />} />
      <Route path="/register" element={<Register />} />
      <Route path="/admin" element={<AdminPanel user={user} />} />
    </Routes>
  );
}

export default AppRouter;

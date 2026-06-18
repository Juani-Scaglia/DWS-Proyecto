import { Routes, Route } from "react-router-dom";

import Home from "../pages/Home";
import EventDetail from "../pages/EventDetail";
import MyTickets from "../pages/MyTickets";
import Login from "../pages/Login";
import Register from "../pages/Register";

function AppRouter() {
  return (
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/event/:id" element={<EventDetail />} />
        <Route path="/tickets" element={<MyTickets />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
      </Routes>
  );
}

export default AppRouter;
import axios from "axios";

const API_URL = "http://localhost:8080/api";

const authHeader = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
});

export const purchaseTickets = (eventId, seatIds) =>
  axios.post(`${API_URL}/tickets/purchase`, { event_id: eventId, seat_ids: seatIds }, authHeader());

export const getMyTickets = () =>
  axios.get(`${API_URL}/tickets/my-tickets`, authHeader());

export const cancelTicket = (ticketId) =>
  axios.post(`${API_URL}/tickets/${ticketId}/cancel`, {}, authHeader());

export const transferTicket = (ticketId, dni) =>
  axios.post(`${API_URL}/tickets/${ticketId}/transfer`, { dni }, authHeader());

export const getEventReport = (eventId) =>
  axios.get(`${API_URL}/admin/events/${eventId}/report`, authHeader());

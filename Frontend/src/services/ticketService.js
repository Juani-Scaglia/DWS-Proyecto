import axios from "axios";

const API_URL = "http://localhost:8080/api";

const authHeader = () => ({
  Authorization: `Bearer ${localStorage.getItem("token")}`,
});

export const getMyTickets = async () => {
  return axios.get(`${API_URL}/tickets/my-tickets`, { headers: authHeader() });
};

export const purchaseTicket = async (eventId) => {
  return axios.post(`${API_URL}/tickets/purchase`, { event_id: eventId }, { headers: authHeader() });
};

export const cancelTicket = async (ticketId) => {
  return axios.post(`${API_URL}/tickets/${ticketId}/cancel`, {}, { headers: authHeader() });
};

export const transferTicket = async (ticketId, dni) => {
  return axios.post(`${API_URL}/tickets/${ticketId}/transfer`, { dni }, { headers: authHeader() });
};

export const getEventReport = (eventId) =>
  axios.get(`${API_URL}/admin/events/${eventId}/report`, { headers: authHeader() });

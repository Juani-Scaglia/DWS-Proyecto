import axios from "axios";

const API_URL = "http://localhost:8080/api";

const authHeader = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
});

export const getEvents = async (category = "") => {
  const params = category ? { category } : {};
  return axios.get(`${API_URL}/events`, { params });
};

export const getEventById = async (id) => {
  return axios.get(`${API_URL}/events/${id}`);
};

export const createEvent = async (data) =>
  axios.post(`${API_URL}/admin/events`, data, authHeader());

export const updateEvent = async (id, data) =>
  axios.put(`${API_URL}/admin/events/${id}`, data, authHeader());

export const deleteEvent = async (id) =>
  axios.delete(`${API_URL}/admin/events/${id}`, authHeader());

export const getEventSeats = async (id) =>
  axios.get(`${API_URL}/events/${id}/seats`);

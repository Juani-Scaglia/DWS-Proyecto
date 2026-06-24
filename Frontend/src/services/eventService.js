import axios from "axios";

const API_URL = "http://localhost:8080/api";

const authHeader = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
});

export const getEvents = (category = "") => {
  const params = category ? { category } : {};
  return axios.get(`${API_URL}/events`, { params });
};

export const getEventById = (id) => axios.get(`${API_URL}/events/${id}`);

export const getEventSeats = (id) => axios.get(`${API_URL}/events/${id}/seats`);

export const createEvent = (data) =>
  axios.post(`${API_URL}/admin/events`, data, authHeader());

export const updateEvent = (id, data) =>
  axios.put(`${API_URL}/admin/events/${id}`, data, authHeader());

export const deleteEvent = (id) =>
  axios.delete(`${API_URL}/admin/events/${id}`, authHeader());

export const uploadImage = (file) => {
  const formData = new FormData();
  formData.append("imagen", file);
  return axios.post(`${API_URL}/admin/upload`, formData, {
    headers: {
      Authorization: `Bearer ${localStorage.getItem("token")}`,
      "Content-Type": "multipart/form-data",
    },
  });
};

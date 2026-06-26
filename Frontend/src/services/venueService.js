import axios from "axios";

const API_URL = "http://localhost:8080/api";

const authHeader = () => ({
  headers: { Authorization: `Bearer ${localStorage.getItem("token")}` },
});

export const getVenues = () => axios.get(`${API_URL}/venues`);

export const getVenueById = (id) => axios.get(`${API_URL}/venues/${id}`);

export const createVenue = (data) =>
  axios.post(`${API_URL}/admin/venues`, data, authHeader());

export const updateVenue = (id, data) =>
  axios.put(`${API_URL}/admin/venues/${id}`, data, authHeader());

export const deleteVenue = (id) =>
  axios.delete(`${API_URL}/admin/venues/${id}`, authHeader());

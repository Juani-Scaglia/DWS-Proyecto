import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const getEvents = async (category = "") => {
  const params = category ? { category } : {};
  return axios.get(`${API_URL}/events`, { params });
};

export const getEventById = async (id) => {
  return axios.get(`${API_URL}/events/${id}`);
};

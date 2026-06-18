import axios from "axios";

const API_URL = "http://localhost:8080";

export const getEvents = async () => {
  return axios.get(`${API_URL}/events`);
};

export const getEventById = async (id) => {
  return axios.get(`${API_URL}/events/${id}`);
};
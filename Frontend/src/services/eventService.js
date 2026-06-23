import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const getEvents = async () => {
  const response = await axios.get(
    `${API_URL}/events`
  );

  return response.data;
};

export const getEventById = async (id) => {
  const response = await axios.get(
    `${API_URL}/events/${id}`
  );

  return response.data;
};
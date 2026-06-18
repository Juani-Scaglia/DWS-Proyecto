import axios from "axios";

const API_URL = "http://localhost:8080";

export const getMyTickets = async (token) => {
  return axios.get(`${API_URL}/tickets/my`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
};
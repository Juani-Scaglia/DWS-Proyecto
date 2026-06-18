import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const login = async (email, password) => {
  return axios.post(`${API_URL}/auth/login`, { email, password });
};

export const register = async (data) => {
  return axios.post(`${API_URL}/auth/register`, data);
};

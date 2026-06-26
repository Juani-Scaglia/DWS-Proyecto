import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const login = (email, password) =>
  axios.post(`${API_URL}/auth/login`, { email, password });

export const register = (data) =>
  axios.post(`${API_URL}/auth/register`, data);

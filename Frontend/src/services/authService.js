/*import axios from "axios";

const API_URL = "http://localhost:8080";

export const login = async (data) => {
  return axios.post(`${API_URL}/login`, data);
};

export const register = async (data) => {
  return axios.post(`${API_URL}/register`, data);
};*/

export const login = async (email, password) => {
  console.log("Login", email, password);
};

export const register = async (user) => {
  console.log("Registro", user);
};

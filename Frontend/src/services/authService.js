import axios from "axios";

const API_URL = "http://localhost:8080/api";

export const login = async (email, password) => {
  const response = await axios.post(
    `${API_URL}/auth/login`,
    {
      email,
      password,
    }
  );

  return response.data;
};

export const register = async (userData) => {
  const response = await axios.post(
    `${API_URL}/auth/register`,
    userData
  );

  return response.data;
};

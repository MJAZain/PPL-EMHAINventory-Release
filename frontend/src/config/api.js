import axios from "axios";

const API_BASE_URL = "http://localhost:8080/api";

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
});

apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 403) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      window.location.href = "/?error=inactive";
    }
    return Promise.reject(error);
  }
);

export const publicClient = axios.create({
  baseURL: API_BASE_URL,
});

export default {
  apiClient,
  publicClient,
};

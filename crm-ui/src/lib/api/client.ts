import axios from "axios";
import { attachInterceptors } from "./interceptors";

export const apiClient = axios.create({
  baseURL: "/api/crm",
  timeout: 15_000,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
  withCredentials: true,
});

attachInterceptors(apiClient);

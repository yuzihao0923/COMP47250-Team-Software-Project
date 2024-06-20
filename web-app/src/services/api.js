// web-app/src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:8080'; // 修改为你的API地址

export const login = (credentials) => axios.post(`${API_URL}/login`, credentials);
export const produce = (message) => axios.post(`${API_URL}/produce`, message);
export const register = (data) => axios.post(`${API_URL}/register`, data);
export const consume = () => axios.get(`${API_URL}/consume`);
export const getLogs = () => axios.get(`${API_URL}/logs`);
export const ackMessage = (id) => axios.post(`${API_URL}/ack`, { id });

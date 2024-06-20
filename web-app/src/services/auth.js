// src/services/authjs
import axios from 'axios';

const API_URL = 'http://localhost:8080';

const login = async (username, password) => {
    const response = await axios.post(`${API_URL}/login`, { username, password });
    localStorage.setItem('token', response.data.token);
};

const logout = () => {
    localStorage.removeItem('token');
};

const getToken = () => {
    return localStorage.getItem('token');
};

export default {
    login,
    logout,
    getToken,
};

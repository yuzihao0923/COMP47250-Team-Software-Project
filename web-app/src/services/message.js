// src/services/message.js
import axios from 'axios';
import authService from './auth';

const API_URL = 'http://localhost:8080';

const sendMessage = async (message) => {
    const token = authService.getToken();
    await axios.post(`${API_URL}/produce`, { message }, {
        headers: { Authorization: `Bearer ${token}` },
    });
};

const consumeMessages = async () => {
    const token = authService.getToken();
    const response = await axios.get(`${API_URL}/consume`, {
        headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
};

const getBrokerLogs = async () => {
    const token = authService.getToken();
    const response = await axios.get(`${API_URL}/logs`, {
        headers: { Authorization: `Bearer ${token}` },
    });
    return response.data;
};

export default {
    sendMessage,
    consumeMessages,
    getBrokerLogs,
};

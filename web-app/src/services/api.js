// src/services/api.js
import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const login = (credentials) => {
  return axios.post(`${API_URL}/login`, credentials).then(response => {
    const { token, role } = response.data;
    if (role === 'broker') {
      localStorage.setItem('token', token); // store broker token
    } else {
      throw new Error('This user is not a broker, please try again');
    }
    return response;
  });
};

export const produce = (message) => {
  const token = localStorage.getItem('token');
  return axios.post(`${API_URL}/produce`, message, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const register = (data) => {
  const token = localStorage.getItem('token');
  return axios.post(`${API_URL}/register`, data, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const consume = () => {
  const token = localStorage.getItem('token');
  return axios.get(`${API_URL}/consume`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const getLogs = () => {
  const token = localStorage.getItem('token');
  return axios.get(`${API_URL}/logs`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const ackMessage = (id) => {
  const token = localStorage.getItem('token');
  return axios.post(`${API_URL}/ack`, { id }, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const connectWebSocket = (role, onMessageCallback) => {
  const token = localStorage.getItem('token');
  let socket = new WebSocket(`ws://localhost:8080/ws?token=${token}`);

  socket.onopen = () => {
    console.log('WebSocket connection established');
  };

  socket.onmessage = (event) => {
    console.log('WebSocket message received:', event.data);
    onMessageCallback(event.data);
  };

  socket.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  socket.onclose = (event) => {
    console.log(`WebSocket connection closed: ${event.code} - ${event.reason}`);
    // Retry connection
    setTimeout(() => {
      console.log('Retrying WebSocket connection...');
      socket = connectWebSocket(role, onMessageCallback);
    }, 3000);
  };

  return socket;
};

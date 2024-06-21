import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const login = (credentials) => {
  return axios.post(`${API_URL}/login`, credentials).then(response => {
    const token = response.data.token;
    localStorage.setItem(`${credentials.role}_token`, token); // store token depends on role
    return response;
  });
};

export const produce = (message, role = 'producer') => {
  const token = localStorage.getItem(`${role}_token`);
  return axios.post(`${API_URL}/produce`, message, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const register = (data, role = 'consumer') => {
  const token = localStorage.getItem(`${role}_token`);
  return axios.post(`${API_URL}/register`, data, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const consume = (role = 'consumer') => {
  const token = localStorage.getItem(`${role}_token`);
  return axios.get(`${API_URL}/consume`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const getLogs = (role) => {
  const token = localStorage.getItem(`${role}_token`);
  return axios.get(`${API_URL}/logs`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const ackMessage = (id, role = 'consumer') => {
  const token = localStorage.getItem(`${role}_token`);
  return axios.post(`${API_URL}/ack`, { id }, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const connectWebSocket = (role, onMessageCallback) => {
  const token = localStorage.getItem(`${role}_token`);
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

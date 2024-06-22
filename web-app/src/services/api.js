import axios from 'axios';

const API_URL = 'http://localhost:8080';

export const login = (credentials) => {
  return axios.post(`${API_URL}/login`, credentials).then(response => {
    const { token, username } = response.data;
    localStorage.setItem(`${username}_token`, token); // store token with username
    return response;
  });
};

// Other API calls using the token
export const produce = (message, username) => {
  const token = localStorage.getItem(`${username}_token`);
  return axios.post(`${API_URL}/produce`, message, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const register = (data, username) => {
  const token = localStorage.getItem(`${username}_token`);
  return axios.post(`${API_URL}/register`, data, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const consume = (username) => {
  const token = localStorage.getItem(`${username}_token`);
  return axios.get(`${API_URL}/consume`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const getLogs = (username) => {
  const token = localStorage.getItem(`${username}_token`);
  return axios.get(`${API_URL}/logs`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const ackMessage = (id, username) => {
  const token = localStorage.getItem(`${username}_token`);
  return axios.post(`${API_URL}/ack`, { id }, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
};

export const connectWebSocket = (username, onMessageCallback) => {
  const token = localStorage.getItem(`${username}_token`);
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
      socket = connectWebSocket(username, onMessageCallback);
    }, 3000);
  };

  return socket;
};

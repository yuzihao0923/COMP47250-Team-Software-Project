// src/services/socket.js
const API_URL = 'ws://localhost:8080/ws';
let socket;

export const connectWebSocket = (username, onMessageCallback) => {
  const token = localStorage.getItem(`${username}_token`);

  if (socket && (socket.readyState === WebSocket.OPEN || socket.readyState === WebSocket.CONNECTING)) {
    return socket;
  }

  socket = new WebSocket(`${API_URL}?token=${token}`);

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
    if (event.code !== 1000) { // If the close event is not normal, retry the connection
      setTimeout(() => {
        console.log('Retrying WebSocket connection...');
        connectWebSocket(username, onMessageCallback);
      }, 3000);
    }
  };

  return socket;
};

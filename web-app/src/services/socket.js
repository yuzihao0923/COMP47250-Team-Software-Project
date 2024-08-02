// src/services/socket.js
const BASE_WS_URL = 'ws://localhost';
const PORTS = [8081, 8082, 8083, 8084];
let sockets = {};

export const connectWebSockets = (user, onMessageCallback) => {
  const token = user.token;

  PORTS.forEach(port => {
    if (!sockets[port] || (sockets[port].readyState !== WebSocket.OPEN && sockets[port].readyState !== WebSocket.CONNECTING)) {
      const wsUrl = `${BASE_WS_URL}:${port}/ws?token=${token}`;
      sockets[port] = new WebSocket(wsUrl);

      sockets[port].onopen = () => {
        console.log(`WebSocket connection established on port ${port}`);
      };

      sockets[port].onmessage = (event) => {
        console.log(`WebSocket message received on port ${port}:`, event.data);
        onMessageCallback(event, port);  // Pass the entire event and port number
      };


      sockets[port].onerror = (error) => {
        console.error(`WebSocket error on port ${port}:`, error);
      };

      sockets[port].onclose = (event) => {
        console.log(`WebSocket connection closed on port ${port}: ${event.code} - ${event.reason}`);
        if (event.code !== 1000) {
          setTimeout(() => {
            console.log(`Retrying WebSocket connection on port ${port}...`);
            connectWebSockets(user, onMessageCallback);
          }, 3000);
        }
      };
    }
  });

  return sockets;
};

export const closeAllSockets = () => {
  Object.values(sockets).forEach(socket => {
    if (socket.readyState === WebSocket.OPEN) {
      socket.close();
    }
  });
};
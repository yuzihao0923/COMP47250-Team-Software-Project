import React, { useEffect, useState } from 'react';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css';
import Head from '../components/Header';

const ConsumerConsole = () => {
  const [logs, setLogs] = useState([]);
  const role = 'consumer';

  useEffect(() => {
    const socket = connectWebSocket(role, (message) => {
      // console.log('Received message:', message);
      if (message.includes('[Consumer]')) {
        setLogs((prevLogs) => [...prevLogs, message]);
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [role]);

  return (
    <div>
      <Head />
    <div className="console-container">
      <h1>Consumer Console</h1>
      <div className="console-logs">
        {logs.map((log, index) => (
          <p key={index} className="console-log">{log}</p>
        ))}
      </div>
    </div>
    </div>

  );
};

export default ConsumerConsole;

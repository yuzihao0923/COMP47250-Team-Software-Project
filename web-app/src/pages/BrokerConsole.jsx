import React, { useEffect, useState } from 'react';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css'

const BrokerConsole = () => {
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    const socket = connectWebSocket((message) => {
      if (message.includes('[Broker]') || message.includes('[Redis]')) {
        setLogs((prevLogs) => [...prevLogs, message]);
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, []);

  return (
    <div className="console-container">
      <h1>Broker Console</h1>
      <div className="console-logs">
        {logs.map((log, index) => (
          <p key={index} className="console-log">{log}</p>
        ))}
      </div>
    </div>
  );
};

export default BrokerConsole;

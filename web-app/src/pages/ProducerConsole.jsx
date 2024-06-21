import React, { useEffect, useState } from 'react';
import { connectWebSocket } from '../services/socket';
import { produce } from '../services/api';
import '../css/Console.css'

const ProducerConsole = () => {
  const [message, setMessage] = useState('');
  const [logs, setLogs] = useState([]);

  useEffect(() => {
    const socket = connectWebSocket((message) => {
      if (message.includes('[Producer]')) {
        setLogs((prevLogs) => [...prevLogs, message]);
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, []);

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await produce({ message });
      console.log('Producer sent message:', response.data);
    } catch (error) {
      console.error('Failed to send message:', error);
    }
  };

  return (
    <div>
      <h1>Producer Console</h1>
      <form onSubmit={handleSubmit}>
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          placeholder="Message"
        />
        <button type="submit">Send Message</button>
      </form>
      <div className="console-container">
        <h1>Producer Console</h1>
        <div className="console-logs">
          {logs.map((log, index) => (
            <p key={index} className="console-log">{log}</p>
          ))}
        </div>
      </div>
    </div>
  );
};

export default ProducerConsole;

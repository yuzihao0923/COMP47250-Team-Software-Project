import React, { useEffect, useState, useRef } from 'react';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css';

const BrokerConsole = () => {
  const [brokerLogs, setBrokerLogs] = useState([]);
  const [producerLogs, setProducerLogs] = useState([]);
  const [consumerLogs, setConsumerLogs] = useState([]);

  const brokerLogsEndRef = useRef(null);
  const producerLogsEndRef = useRef(null);
  const consumerLogsEndRef = useRef(null);

  useEffect(() => {
    const socket = connectWebSocket((message) => {
      if (message.includes('[Broker]') || message.includes('[Redis]')) {
        setBrokerLogs((prevLogs) => [...prevLogs, message]);
      } else if (message.includes('[Producer]')) {
        setProducerLogs((prevLogs) => [...prevLogs, message]);
      } else if (message.includes('[Consumer]')) {
        setConsumerLogs((prevLogs) => [...prevLogs, message]);
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, []);

  useEffect(() => {
    brokerLogsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [brokerLogs]);

  useEffect(() => {
    producerLogsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [producerLogs]);

  useEffect(() => {
    consumerLogsEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [consumerLogs]);

  return (
    <div className="console-container">
      <h1>Broker Console</h1>
      <div className="log-section broker-logs">
        <h2>Broker & Redis Logs</h2>
        <div className="console-logs">
          {brokerLogs.map((log, index) => (
            <p key={index} className="console-log">{log}</p>
          ))}
          <div ref={brokerLogsEndRef} />
        </div>
      </div>
      <div className="log-section producer-consumer-logs">
        <div className="log-subsection producer-logs">
          <h2>Producer Logs</h2>
          <div className="console-logs">
            {producerLogs.map((log, index) => (
              <p key={index} className="console-log">{log}</p>
            ))}
            <div ref={producerLogsEndRef} />
          </div>
        </div>
        <div className="log-subsection consumer-logs">
          <h2>Consumer Logs</h2>
          <div className="console-logs">
            {consumerLogs.map((log, index) => (
              <p key={index} className="console-log">{log}</p>
            ))}
            <div ref={consumerLogsEndRef} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default BrokerConsole;

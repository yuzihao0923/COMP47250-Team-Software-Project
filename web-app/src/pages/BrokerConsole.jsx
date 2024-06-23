import React, { useEffect, useState, useRef } from 'react';
import { useSelector } from 'react-redux';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css';

const BrokerConsole = () => {
  const [brokerLogs, setBrokerLogs] = useState([]);
  const [producerLogs, setProducerLogs] = useState([]);
  const [consumerLogs, setConsumerLogs] = useState([]);
  const [totalProducerMessages, setTotalProducerMessages] = useState(0);
  const [intervalMessageCount, setIntervalMessageCount] = useState(0);
  const [chartData, setChartData] = useState({
    labels: [],  // Time labels for the chart
    datasets: [
      {
        label: 'Messages per Interval',
        data: [],  // Data points for the chart
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.1,
        fill: false,
      }
    ]
  });

  const user = useSelector(state => state.user);
  const updateInterval = useRef(null);

  useEffect(() => {
    const socket = connectWebSocket(user.username, (message) => {
      const cleanedMessage = message.replace(/"/g, '');
      if (cleanedMessage.includes('[Producer')) {
        setProducerLogs(prevLogs => [...prevLogs, cleanedMessage]);
        setTotalProducerMessages(prevCount => prevCount + 1);
        setIntervalMessageCount(prevCount => prevCount + 1);
      } else if (cleanedMessage.includes('[Broker') || cleanedMessage.includes('[Redis')) {
        setBrokerLogs(prevLogs => [...prevLogs, cleanedMessage]);
      } else if (cleanedMessage.includes('[Consumer')) {
        setConsumerLogs(prevLogs => [...prevLogs, cleanedMessage]);
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [user.username]);

  useEffect(() => {
    updateInterval.current = setInterval(() => {
      setChartData(prevData => {
        const currentTime = new Date().toLocaleTimeString();
        const newLabels = [...prevData.labels, currentTime];
        const newDataPoints = [...prevData.datasets[0].data, intervalMessageCount];

        if (newLabels.length > 60) {
          newLabels.shift();
          newDataPoints.shift();
        }

        return {
          ...prevData,
          labels: newLabels,
          datasets: [{ ...prevData.datasets[0], data: newDataPoints }]
        };
      });
      setIntervalMessageCount(0); // Reset interval count after updating the chart
    }, 5000);

    return () => clearInterval(updateInterval.current);
  }, [intervalMessageCount]);

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
        y: {
            beginAtZero: true, // Start the scale at zero
            ticks: {
                stepSize: 1, // Ensures the scale steps at integer values
                callback: function(value) {
                    if (value % 1 === 0) { // Display only integers
                        return value;
                    }
                }
            }
        }
    },
    animation: {
        duration: 500
    }
};



  return (
    <div className="console-container">
      <h1>Broker Console</h1>
      <div>Total Producer Messages Sent: {totalProducerMessages}</div>
      {/* Adjust the chart size by modifying the maxWidth, width, and height properties */}
      <div style={{ width: '95%', maxWidth: '1500px', height: '500px', margin: '0 auto' }}> 
        <Line data={chartData} options={chartOptions} />
      </div>
      <div className="log-section broker-logs">
        <h2>Broker & Redis Logs</h2>
        <div className="console-logs">
          {brokerLogs.map((log, index) => (
            <p key={index} className="console-log">{log}</p>
          ))}
        </div>
      </div>
      <div className="log-section producer-consumer-logs">
        <div className="log-subsection producer-logs">
          <h2>Producer Logs</h2>
          <div className="console-logs">
            {producerLogs.map((log, index) => (
              <p key={index} className="console-log">{log}</p>
            ))}
          </div>
        </div>
        <div className="log-subsection consumer-logs">
          <h2>Consumer Logs</h2>
          <div className="console-logs">
            {consumerLogs.map((log, index) => (
              <p key={index} className="console-log">{log}</p>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default BrokerConsole;

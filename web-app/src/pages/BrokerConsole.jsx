import React, { useEffect, useState, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import 'chartjs-adapter-date-fns';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css';
import Card from '../components/Card';
import { SendOutlined, ShareAltOutlined, ApiOutlined } from '@ant-design/icons';
import Logs from '../components/Logs';
import { addProducerLog, resetProducerIntervalCounts, updateProducerChartData } from '../store/producerSlice'
import { addConsumerLog, resetConsumerIntervalCounts, updateConsumerChartData } from '../store/consumerSlice'
import { addBrokerLog, resetBrokerIntervalCounts, addBrokerAcknowledgedMessage, updateBrokerChartData } from '../store/brokerSlice'

const BrokerConsole = () => {
  // const [brokerLogs, setBrokerLogs] = useState([]);
  // const [producerLogs, setProducerLogs] = useState([]);
  // const [consumerLogs, setConsumerLogs] = useState([]);
  // const [totalProducerMessages, setTotalProducerMessages] = useState(0);
  // const [totalConsumerReceivedMessages, setTotalConsumerReceivedMessages] = useState(0);
  // const [totalBrokerAcknowledgedMessages, setTotalBrokerAcknowledgedMessages] = useState(0);
  // const [intervalProducerMessageCount, setIntervalProducerMessageCount] = useState(0);
  // const [intervalConsumerReceivedMessageCount, setIntervalConsumerReceivedMessageCount] = useState(0);
  // const [intervalBrokerAcknowledgedMessageCount, setIntervalBrokerAcknowledgedMessageCount] = useState(0);

  // const [producerChartData, setProducerChartData] = useState({
  //   labels: [],
  //   datasets: [
  //     {
  //       label: 'Producer Messages per Interval',
  //       data: [],
  //       borderColor: 'rgb(75, 192, 192)',
  //       tension: 0.1,
  //       fill: false,
  //     }
  //   ]
  // });
  // const [consumerChartData, setConsumerChartData] = useState({
  //   labels: [],
  //   datasets: [
  //     {
  //       label: 'Consumer Messages Received per Interval',
  //       data: [],
  //       borderColor: 'rgb(255, 99, 132)',
  //       tension: 0.1,
  //       fill: false,
  //     }
  //   ]
  // });
  // const [brokerChartData, setBrokerChartData] = useState({
  //   labels: [],
  //   datasets: [
  //     {
  //       label: 'Broker Acknowledged Messages per Interval',
  //       data: [],
  //       borderColor: 'rgb(54, 162, 235)',
  //       tension: 0.1,
  //       fill: false,
  //     }
  //   ]
  // });

  const dispatch = useDispatch();
  const brokerLogs = useSelector((state) => state.broker.brokerLogs);
  const producerLogs = useSelector((state) => state.producer.producerLogs);
  const consumerLogs = useSelector((state) => state.consumer.consumerLogs);
  const totalProducerMessages = useSelector((state) => state.producer.totalProducerMessages);
  const totalConsumerReceivedMessages = useSelector((state) => state.consumer.totalConsumerReceivedMessages);
  const totalBrokerAcknowledgedMessages = useSelector((state) => state.broker.totalBrokerAcknowledgedMessages);
  const intervalProducerMessageCount = useSelector((state) => state.producer.intervalProducerMessageCount);
  const intervalConsumerReceivedMessageCount = useSelector((state) => state.consumer.intervalConsumerReceivedMessageCount);
  const intervalBrokerAcknowledgedMessageCount = useSelector((state) => state.broker.intervalBrokerAcknowledgedMessageCount);
  const producerChartData = useSelector((state) => state.producer.producerChartData);
  const consumerChartData = useSelector((state) => state.consumer.consumerChartData);
  const brokerChartData = useSelector((state) => state.broker.brokerChartData);

  const user = useSelector(state => state.user);
  const updateInterval = useRef(null);

  // useEffect(() => {
  //   const socket = connectWebSocket(user, (message) => {
  //     const cleanedMessage = message.replace(/"/g, '');
  //     if (cleanedMessage.includes('[Producer')) {
  //       setProducerLogs(prevLogs => [...prevLogs, cleanedMessage]);
  //       setTotalProducerMessages(prevCount => prevCount + 1);
  //       setIntervalProducerMessageCount(prevCount => prevCount + 1);
  //     } else if (cleanedMessage.includes('[Broker') || cleanedMessage.includes('[Redis')) {
  //       setBrokerLogs(prevLogs => [...prevLogs, cleanedMessage]);
  //       if (cleanedMessage.includes('acknowledged successfully')) {
  //         setTotalBrokerAcknowledgedMessages(prevCount => prevCount + 1);
  //         setIntervalBrokerAcknowledgedMessageCount(prevCount => prevCount + 1);
  //       }
  //     } else if (cleanedMessage.includes('[Consumer')) {
  //       setConsumerLogs(prevLogs => [...prevLogs, cleanedMessage]);
  //       if (cleanedMessage.includes('received')) {
  //         setTotalConsumerReceivedMessages(prevCount => prevCount + 1);
  //         setIntervalConsumerReceivedMessageCount(prevCount => prevCount + 1);
  //       }
  //     }
  //   });

  //   return () => {
  //     if (socket && socket.readyState === WebSocket.OPEN) {
  //       socket.close();
  //     }
  //   };
  // }, [user]);

  useEffect(() => {
    const socket = connectWebSocket(user, (message) => {
      const cleanedMessage = message.replace(/"/g, '');
      if (cleanedMessage.includes('[Producer')) {
        dispatch(addProducerLog(cleanedMessage));
      } else if (cleanedMessage.includes('[Broker') || cleanedMessage.includes('[Redis')) {
        dispatch(addBrokerLog(cleanedMessage));
        if (cleanedMessage.includes('acknowledged successfully')) {
          dispatch(addBrokerAcknowledgedMessage());
        }
      } else if (cleanedMessage.includes('[Consumer')) {
        dispatch(addConsumerLog(cleanedMessage));
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [dispatch, user]);

  // useEffect(() => {
  //   updateInterval.current = setInterval(() => {
  //     const currentTime = new Date();

  //     setProducerChartData(prevData => {
  //       const newLabels = [...prevData.labels, currentTime];
  //       const newDataPoints = [...prevData.datasets[0].data, intervalProducerMessageCount];

  //       if (newLabels.length > 20) {
  //         newLabels.shift();
  //         newDataPoints.shift();
  //       }

  //       return {
  //         ...prevData,
  //         labels: newLabels,
  //         datasets: [{ ...prevData.datasets[0], data: newDataPoints }]
  //       };
  //     });

  //     setConsumerChartData(prevData => {
  //       const newLabels = [...prevData.labels, currentTime];
  //       const newReceivedDataPoints = [...prevData.datasets[0].data, intervalConsumerReceivedMessageCount];

  //       if (newLabels.length > 20) {
  //         newLabels.shift();
  //         newReceivedDataPoints.shift();
  //       }

  //       return {
  //         ...prevData,
  //         labels: newLabels,
  //         datasets: [{ ...prevData.datasets[0], data: newReceivedDataPoints }]
  //       };
  //     });

  //     setBrokerChartData(prevData => {
  //       const newLabels = [...prevData.labels, currentTime];
  //       const newAcknowledgedDataPoints = [...prevData.datasets[0].data, intervalBrokerAcknowledgedMessageCount];

  //       if (newLabels.length > 20) {
  //         newLabels.shift();
  //         newAcknowledgedDataPoints.shift();
  //       }

  //       return {
  //         ...prevData,
  //         labels: newLabels,
  //         datasets: [{ ...prevData.datasets[0], data: newAcknowledgedDataPoints }]
  //       };
  //     });

  //     setIntervalProducerMessageCount(0);
  //     setIntervalConsumerReceivedMessageCount(0);
  //     setIntervalBrokerAcknowledgedMessageCount(0);
  //   }, 3000);

  //   return () => clearInterval(updateInterval.current);
  // }, [intervalProducerMessageCount, intervalConsumerReceivedMessageCount, intervalBrokerAcknowledgedMessageCount]);

  useEffect(() => {
    updateInterval.current = setInterval(() => {
      const currentTime = new Date();

      dispatch(updateProducerChartData({
        labels: [...producerChartData.labels, currentTime],
        data: [...producerChartData.datasets[0].data, intervalProducerMessageCount]
      }));

      dispatch(updateConsumerChartData({
        labels: [...consumerChartData.labels, currentTime],
        data: [...consumerChartData.datasets[0].data, intervalConsumerReceivedMessageCount]
      }));

      dispatch(updateBrokerChartData({
        labels: [...brokerChartData.labels, currentTime],
        data: [...brokerChartData.datasets[0].data, intervalBrokerAcknowledgedMessageCount]
      }));

      dispatch(resetProducerIntervalCounts());
      dispatch(resetConsumerIntervalCounts());
      dispatch(resetBrokerIntervalCounts());
    }, 3000);

    return () => clearInterval(updateInterval.current);
  }, [dispatch, intervalProducerMessageCount, intervalConsumerReceivedMessageCount, intervalBrokerAcknowledgedMessageCount, producerChartData, consumerChartData, brokerChartData]);

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    scales: {
      x: {
        type: 'time',
        time: {
          unit: 'second',
          stepSize: 10,
          displayFormats: {
            second: 'h:mm:ss a',
          }
        },
        ticks: {
          maxRotation: 0,
          minRotation: 0,
          source: 'auto',
          autoSkip: true,
          maxTicksLimit: 20,
          callback: function (value) {
            const date = new Date(value);
            return date.toLocaleTimeString();
          }
        },
        min: new Date(Date.now() - 60000),
        max: new Date()
      },
      y: {
        beginAtZero: true,
        ticks: {
          stepSize: 1,
          callback: function (value) {
            if (value % 1 === 0) {
              return value;
            }
          }
        }
      }
    },
    animation: {
      duration: 300
    }
  };

  return (
    <div className="console-container">
      <h1>Broker Console</h1>
      <div className='card-area mb-10'>
        <Card logoBackground='bg-sky-200' logo={<SendOutlined />} data={totalProducerMessages} dataTitle='Total Producer Messages Sent' />
        <Card logoBackground='bg-sky-300' logo={<ApiOutlined />} data={totalBrokerAcknowledgedMessages} dataTitle='Total Broker Acknowledged Messages' />
        <Card logoBackground='bg-sky-400' logo={<ShareAltOutlined />} data={totalConsumerReceivedMessages} dataTitle='Total Consumer Messages Received' />
      </div>

      <h1>Monitor Chart</h1>
      <div className="charts-container">
        <div className="chart-wrapper">
          <Line data={producerChartData} options={chartOptions} />
        </div>
        <div className="chart-wrapper">
          <Line data={consumerChartData} options={chartOptions} />
        </div>
        <div className="chart-wrapper">
          <Line data={brokerChartData} options={chartOptions} />
        </div>
      </div>

      <h1>Components Logs</h1>
      <div className='w-full'>
        {/* Broker & Redis Logs */}
        <Logs logsTitle='Broker & Redis Logs' logsBackgroundColor='bg-orange-100' logsData={brokerLogs} />
        {/* Producer Logs */}
        <Logs logsTitle='Producer Logs' logsBackgroundColor='bg-indigo-200' logsData={producerLogs} />
        {/* Consumer Logs */}
        <Logs logsTitle='Consumer Logs' logsBackgroundColor='bg-cyan-100' logsData={consumerLogs} />
      </div>
    </div>
  );
};

export default BrokerConsole;

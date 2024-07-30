import React, { useEffect, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { Line } from 'react-chartjs-2';
import 'chart.js/auto';
import 'chartjs-adapter-date-fns';
import { connectWebSocket } from '../services/socket';
import '../css/Console.css';
import Card from '../components/Card';
// import { SendOutlined, ShareAltOutlined, ApiOutlined } from '@ant-design/icons';
import { SendOutlined} from '@ant-design/icons';
import Logs from '../components/Logs';
import { addProducerLog, resetProducerIntervalCounts, updateProducerChartData } from '../store/producerSlice'
import { addConsumerLog, resetConsumerIntervalCounts, updateConsumerChartData } from '../store/consumerSlice'
import { addBrokerLog, resetBrokerIntervalCounts, addBrokerAcknowledgedMessage, updateBrokerChartData } from '../store/brokerSlice'
import { addProxyLog} from '../store/proxySlice';

const BrokerConsole = () => {
  const dispatch = useDispatch();
  const brokerLogs = useSelector((state) => state.broker.brokerLogs);
  const producerLogs = useSelector((state) => state.producer.producerLogs);
  const consumerLogs = useSelector((state) => state.consumer.consumerLogs);
  const proxyLogs = useSelector((state) => state.proxy.proxyLogs);
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

  useEffect(() => {
    const socket = connectWebSocket(user, (message) => {
      console.log("Received message:", message);
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
      } else if (cleanedMessage.includes('[ProxyServer')) {
        dispatch(addProxyLog(cleanedMessage));
      }
    });

    return () => {
      if (socket && socket.readyState === WebSocket.OPEN) {
        socket.close();
      }
    };
  }, [dispatch, user]);


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
      <h1 className="hollow-fluorescent-edge">Broker Console</h1>
      <div className='card-area mb-10'>
      <Card className="card producer-card" logoBackground='bg-dark-200' logo={<SendOutlined />} data={totalProducerMessages} dataTitle='Total Producer Messages Sent'  />
      <Card className="card broker-card" logoBackground='bg-dark-300' logo={<SendOutlined />} data={totalBrokerAcknowledgedMessages} dataTitle='Total Broker Acknowledged Messages' />
      <Card className="card consumer-card" logoBackground='bg-dark-400' logo={<SendOutlined />} data={totalConsumerReceivedMessages} dataTitle='Total Consumer Messages Received' />
      </div>

      <h1 className="hollow-fluorescent-edge">Monitor Chart</h1>
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

      <h1 className="hollow-fluorescent-edge">Components Logs</h1>
      <div className='w-full' style={{ backgroundColor: '#121212' }}> {}
        <Logs logsTitle='Proxy Logs' logsData={proxyLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        {/* Broker & Redis Logs */}
        <Logs logsTitle='Broker & Redis Logs' logsData={brokerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        {/* Producer Logs */}
        <Logs logsTitle='Producer Logs' logsData={producerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        {/* Consumer Logs */}
        <Logs logsTitle='Consumer Logs' logsData={consumerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
      </div>
    </div>
  );
};

export default BrokerConsole;

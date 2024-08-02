import React, { useEffect, useState, useCallback } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { selectAllLogs, selectAllMetrics } from '../store/selectors';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';
import { connectWebSockets, closeAllSockets } from '../services/socket';
import '../css/Console.css';
import Card from '../components/Card';
import Logs from '../components/Logs';
import { SendOutlined } from '@ant-design/icons';
import { batchAddLogs } from '../store/logSlice';
import { incrementProducerMessage, resetProducerIntervalCounts } from '../store/producerMetrics'
import { incrementConsumerMessage, resetConsumerIntervalCounts } from '../store/consumerMetrics'
import { addBrokerAcknowledgedMessage, resetBrokerIntervalCounts } from '../store/brokerMetrics'
import { incrementProxyMessage, resetProxyIntervalCounts } from '../store/proxyMetrics'


const BrokerConsole = () => {
  const dispatch = useDispatch();
  const logs = useSelector(selectAllLogs);
  const metrics = useSelector(selectAllMetrics);
  const user = useSelector(state => state.user);

  const [chartData, setChartData] = useState([]);

  const processBatchedMessages = useCallback((messages) => {
    console.log('Raw messages:', messages);

    // Remove surrounding quotes and split by newline
    const cleanedMessage = messages.replace(/^"|"$/g, '');
    const logEntries = cleanedMessage.split('\\n').filter(entry => entry.trim() !== '');

    console.log('Processed log entries:', logEntries);

    dispatch(batchAddLogs(logEntries));

    let brokerCount = 0, consumerCount = 0, producerCount = 0, proxyCount = 0;

    logEntries.forEach(message => {
      if (message.includes('[Producer')) {
        producerCount++;
      } else if (message.includes('[Consumer')) {
        consumerCount++;
      } else if (message.includes('[Broker') || message.includes('[Redis')) {
        if (message.includes('acknowledged successfully')) {
          brokerCount++;
        }
      } else if (message.includes('[ProxyServer')) {
        proxyCount++;
      }
    });

    dispatch(incrementProducerMessage(producerCount));
    dispatch(incrementConsumerMessage(consumerCount));
    dispatch(addBrokerAcknowledgedMessage(brokerCount));
    dispatch(incrementProxyMessage(proxyCount));
  }, [dispatch]);

  useEffect(() => {
    const handleWebSocketMessage = (event, port) => {
      console.log(`WebSocket message received on port ${port}:`, event.data);

      if (typeof event.data === 'string') {
        processBatchedMessages(event.data);
      } else {
        console.error('Unexpected WebSocket message format:', event.data);
      }
    };

    const _sockets = connectWebSockets(user, handleWebSocketMessage);

    return () => {
      closeAllSockets();
    };
  }, [user, processBatchedMessages]);

  useEffect(() => {
    const initializeChartData = () => {
      const now = new Date();
      const oneMinuteAgo = new Date(now.getTime() - 60000);
      const initialData = Array.from({ length: 21 }, (_, index) => {
        const time = new Date(oneMinuteAgo.getTime() + index * 3000);
        return {
          time: time.toISOString(),
          broker: 0,
          consumer: 0,
          producer: 0,
        };
      });
      setChartData(initialData);
    };

    initializeChartData();

    const chartUpdateInterval = setInterval(() => {
      setChartData(prevData => {
        const newDataPoint = {
          time: new Date().toISOString(),
          broker: metrics.intervalBrokerAcknowledgedMessageCount,
          consumer: metrics.intervalConsumerReceivedMessageCount,
          producer: metrics.intervalProducerMessageCount,
        };
        const newData = [...prevData.slice(1), newDataPoint];

        // Dispatch actions to reset interval counts after updating the chart
        dispatch(resetBrokerIntervalCounts());
        dispatch(resetConsumerIntervalCounts());
        dispatch(resetProducerIntervalCounts());
        dispatch(resetProxyIntervalCounts());

        return newData;
      });
    }, 3000);

    return () => clearInterval(chartUpdateInterval);
  }, [dispatch]);

  // Update chart data when metrics change
  useEffect(() => {
    setChartData(prevData => {
      const lastIndex = prevData.length - 1;
      const updatedLastPoint = {
        ...prevData[lastIndex],
        broker: metrics.intervalBrokerAcknowledgedMessageCount,
        consumer: metrics.intervalConsumerReceivedMessageCount,
        producer: metrics.intervalProducerMessageCount,
      };
      return [...prevData.slice(0, lastIndex), updatedLastPoint];
    });
  }, [metrics]);

  const formatXAxis = (tickItem) => {
    return new Date(tickItem).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
  };

  return (
    <div className="console-container">
      <h1 className="hollow-fluorescent-edge">Broker Console</h1>
      <div className='card-area mb-10'>
        <Card
          className="card broker-card"
          logoBackground='bg-dark-200'
          logo={<SendOutlined />}
          data={metrics.totalBrokerAcknowledgedMessages}
          dataTitle='Total Broker Acknowledged Messages'
        />
        <Card
          className="card consumer-card"
          logoBackground='bg-dark-300'
          logo={<SendOutlined />}
          data={metrics.totalConsumerReceivedMessages}
          dataTitle='Total Consumer Messages Received'
        />
        <Card
          className="card producer-card"
          logoBackground='bg-dark-400'
          logo={<SendOutlined />}
          data={metrics.totalProducerMessages}
          dataTitle='Total Producer Messages Sent'
        />
      </div>

      <h1 className="hollow-fluorescent-edge">Monitor Chart</h1>
      <div className="charts-container">
        <ResponsiveContainer width="100%" height={500}>
          <LineChart data={chartData}>
            <CartesianGrid horizontal={true} vertical={false} />
            <XAxis
              dataKey="time"
              tickFormatter={formatXAxis}
              interval="preserveEnd"
              minTickGap={50}
              ticks={[chartData[0]?.time, chartData[20]?.time]}
            />
            <YAxis />
            <Tooltip labelFormatter={formatXAxis} />
            <Legend />
            <Line
              type="stepAfter"
              dataKey="broker"
              stroke="#8884d8"
              name="Broker Messages"
              isAnimationActive={false}
            />
            <Line
              type="stepAfter"
              dataKey="consumer"
              stroke="#82ca9d"
              name="Consumer Messages"
              isAnimationActive={false}
            />
            <Line
              type="stepAfter"
              dataKey="producer"
              stroke="#ffc658"
              name="Producer Messages"
              isAnimationActive={false}
            />
          </LineChart>
        </ResponsiveContainer>
      </div>

      <h1 className="hollow-fluorescent-edge">Components Logs</h1>
      <div className='w-full' style={{ backgroundColor: '#121212' }}>
        <Logs logsTitle='Broker & Redis Logs' logsData={logs.brokerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        <Logs logsTitle='Consumer Logs' logsData={logs.consumerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        <Logs logsTitle='Producer Logs' logsData={logs.producerLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
        <Logs logsTitle='Proxy Logs' logsData={logs.proxyLogs} style={{ backgroundColor: 'transparent', boxShadow: '0 4px 6px rgba(0, 0, 0, 0.5)' }} />
      </div>
    </div>
  );
};

export default BrokerConsole;
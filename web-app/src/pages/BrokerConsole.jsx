import React, { useEffect, useState, useCallback, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { selectAllLogs, selectAllMetrics } from '../store/selectors';
import ReactECharts from 'echarts-for-react';
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

  const [chartOption, setChartOption] = useState({});
  const messageCountsRef = useRef({ broker: 0, consumer: 0, producer: 0 });
  const chartDataRef = useRef([]);
  const chartRef = useRef(null);

  const processBatchedMessages = useCallback((messages) => {
    console.log('Raw messages:', messages);

    const cleanedMessage = messages.replace(/^"|"$/g, '');
    const logEntries = cleanedMessage.split('\\n').filter(entry => entry.trim() !== '');

    console.log('Processed log entries:', logEntries);

    dispatch(batchAddLogs(logEntries));

    let totalProcessed = { broker: 0, consumer: 0, producer: 0 };

    logEntries.forEach(message => {
      if (message.includes('[Producer')) {
        totalProcessed.producer++;
        dispatch(incrementProducerMessage());
      } else if (message.includes('] received')) {
        totalProcessed.consumer++;
        dispatch(incrementConsumerMessage());
      } else if (message.includes('[Broker') || message.includes('[Redis')) {
        if (message.includes('acknowledged successfully')) {
          totalProcessed.broker++;
          dispatch(addBrokerAcknowledgedMessage());
        }
      } else if (message.includes('[ProxyServer')) {
        dispatch(incrementProxyMessage());
      }
    });

    messageCountsRef.current.broker += totalProcessed.broker;
    messageCountsRef.current.consumer += totalProcessed.consumer;
    messageCountsRef.current.producer += totalProcessed.producer;

    console.log('Total messages processed:', totalProcessed);
    console.log('Accumulated message counts:', messageCountsRef.current);
    return totalProcessed;
  }, [dispatch]);

  const updateChartData = useCallback(() => {
    const now = Date.now();
    const newDataPoint = [
      now,
      messageCountsRef.current.broker,
      messageCountsRef.current.consumer,
      messageCountsRef.current.producer
    ];

    console.log('New data point:', newDataPoint);

    // Update chartDataRef
    chartDataRef.current = [...chartDataRef.current.slice(-59), newDataPoint];

    // Update chart option
    setChartOption(prevOption => {
      const updatedOption = {
        ...prevOption,
        series: [
          {
            name: 'Broker',
            data: chartDataRef.current.map(item => [item[0], item[1]])
          },
          {
            name: 'Consumer',
            data: chartDataRef.current.map(item => [item[0], item[2]])
          },
          {
            name: 'Producer',
            data: chartDataRef.current.map(item => [item[0], item[3]])
          }
        ]
      };

      // Force chart update
      if (chartRef.current) {
        chartRef.current.getEchartsInstance().setOption(updatedOption);
      }

      return updatedOption;
    });

    // Reset message counts for the next interval
    messageCountsRef.current = { broker: 0, consumer: 0, producer: 0 };

    // Reset interval counts in Redux store
    dispatch(resetBrokerIntervalCounts());
    dispatch(resetConsumerIntervalCounts());
    dispatch(resetProducerIntervalCounts());
    dispatch(resetProxyIntervalCounts());
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

    // Set up interval for regular chart updates
    const intervalId = setInterval(updateChartData, 1000);

    return () => {
      closeAllSockets();
      clearInterval(intervalId);
    };
  }, [user, processBatchedMessages, updateChartData]);

  useEffect(() => {
    // Initialize chart data and option
    const now = Date.now();
    chartDataRef.current = Array.from({ length: 60 }, (_, index) => [
      now - (59 - index) * 1000,
      0,
      0,
      0
    ]);

    setChartOption({
      title: {
        text: 'Message Rates (Last 60 Seconds)',
        textStyle: { color: '#ffffff' }
      },
      tooltip: {
        trigger: 'axis',
        formatter: function (params) {
          const time = new Date(params[0].value[0]).toLocaleTimeString();
          return params.reduce((acc, param) => {
            return acc + `${param.seriesName}: ${param.value[1]}/s<br/>`;
          }, `${time}<br/>`);
        }
      },
      legend: {
        data: ['Broker', 'Consumer', 'Producer'],
        textStyle: { color: '#ffffff' }
      },
      xAxis: {
        type: 'time',
        splitLine: { show: false },
        axisLabel: { color: '#ffffff' },
        min: 'dataMin',
        max: 'dataMax'
      },
      yAxis: {
        type: 'value',
        splitLine: { show: false },
        axisLabel: { color: '#ffffff' }
      },
      series: [
        {
          name: 'Broker',
          type: 'line',
          showSymbol: false,
          data: [],
          smooth: 0.2
        },
        {
          name: 'Consumer',
          type: 'line',
          showSymbol: false,
          data: [],
          smooth: 0.2
        },
        {
          name: 'Producer',
          type: 'line',
          showSymbol: false,
          data: [],
          smooth: 0.2
        }
      ],
      backgroundColor: 'transparent'
    });
  }, []);
  
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
      <div className="charts-container" style={{ height: '500px' }}>
        <ReactECharts
          ref={chartRef}
          option={chartOption}
          notMerge={false}
          lazyUpdate={true}
          style={{ height: '100%' }}
        />
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
import { createSelector } from 'reselect';

const selectProducerMetrics = state => state.producerMetrics;
const selectConsumerMetrics = state => state.consumerMetrics;
const selectBrokerMetrics = state => state.brokerMetrics;

export const selectAllLogs = state => state.logs;

export const selectAllMetrics = createSelector(
    [selectProducerMetrics, selectConsumerMetrics, selectBrokerMetrics],
    (producer, consumer, broker) => ({
        totalProducerMessages: producer.totalProducerMessages,
        totalConsumerReceivedMessages: consumer.totalConsumerReceivedMessages,
        totalBrokerAcknowledgedMessages: broker.totalBrokerAcknowledgedMessages,
        intervalProducerMessageCount: producer.intervalProducerMessageCount,
        intervalConsumerReceivedMessageCount: consumer.intervalConsumerReceivedMessageCount,
        intervalBrokerAcknowledgedMessageCount: broker.intervalBrokerAcknowledgedMessageCount,
    })
);
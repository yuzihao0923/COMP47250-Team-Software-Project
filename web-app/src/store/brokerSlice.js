import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    brokerLogs: [],
    totalBrokerAcknowledgedMessages: 0,
    intervalBrokerAcknowledgedMessageCount: 0,
    brokerChartData: {
        labels: [],
        datasets: [
            {
                label: 'Broker Acknowledged Messages per Interval',
                data: [],
                borderColor: 'rgb(54, 162, 235)',
                tension: 0.1,
                fill: false,
            },
        ],
    },
};

const brokerSlice = createSlice({
    name: 'broker',
    initialState: initialState,
    reducers: {
        addBrokerLog: (state, action) => {
            state.brokerLogs.push(action.payload);
        },
        resetIntervalCounts: (state) => {
            state.intervalBrokerAcknowledgedMessageCount = 0;
        },
        addBrokerAcknowledgedMessage: (state) => {
            state.totalBrokerAcknowledgedMessages += 1;
            state.intervalBrokerAcknowledgedMessageCount += 1;
        },
        updateBrokerChartData: (state, action) => {
            state.brokerChartData.labels = action.payload.labels;
            state.brokerChartData.datasets[0].data = action.payload.data;
        },
    },
});

export const { addConsumerLog, resetIntervalCounts, updateConsumerChartData } = brokerSlice.actions;
export default brokerSlice.reducer;

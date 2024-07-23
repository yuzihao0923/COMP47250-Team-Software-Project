import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    consumerLogs: [],
    totalConsumerReceivedMessages: 0,
    intervalConsumerReceivedMessageCount: 0,
    consumerChartData: {
        labels: [],
        datasets: [
            {
                label: 'Consumer Messages Received per Interval',
                data: [],
                borderColor: 'rgb(255, 99, 132)',
                tension: 0.1,
                fill: false,
            },
        ],
    },
};

const consumerSlice = createSlice({
    name: 'consumer',
    initialState: initialState,
    reducers: {
        addConsumerLog: (state, action) => {
            state.consumerLogs.push(action.payload);
            if (action.payload.includes('received')) {
                state.totalConsumerReceivedMessages += 1;
                state.intervalConsumerReceivedMessageCount += 1;
            }
        },
        resetIntervalCounts: (state) => {
            state.intervalConsumerReceivedMessageCount = 0;
        },
        updateConsumerChartData: (state, action) => {
            state.consumerChartData.labels = action.payload.labels;
            state.consumerChartData.datasets[0].data = action.payload.data;
        },
    },
});

export const { addConsumerLog, resetIntervalCounts, updateConsumerChartData } = consumerSlice.actions;
export default consumerSlice.reducer;

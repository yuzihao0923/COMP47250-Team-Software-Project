import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    producerLogs: [],
    totalProducerMessages: 0,
    intervalProducerMessageCount: 0,
    producerChartData: {
        labels: [],
        datasets: [
            {
                label: 'Producer Messages per Interval',
                data: [],
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1,
                fill: false,
            },
        ],
    }
};

const producerSlice = createSlice({
    name: 'producer',
    initialState: initialState,
    reducers: {
        addProducerLog: (state, action) => {
            state.producerLogs.push(action.payload);
            state.totalProducerMessages += 1;
            state.intervalProducerMessageCount += 1;
        },
        resetProducerIntervalCounts: (state) => {
            state.intervalProducerMessageCount = 0;
        },
        updateProducerChartData: (state, action) => {
            state.producerChartData.labels = action.payload.labels;
            state.producerChartData.datasets[0].data = action.payload.data;
        },
    },
});

export const { addProducerLog, resetProducerIntervalCounts, updateProducerChartData } = producerSlice.actions;
export default producerSlice.reducer;

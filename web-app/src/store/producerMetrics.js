import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalProducerMessages: 0,
    intervalProducerMessageCount: 0,
};

const producerMetricsSlice = createSlice({
    name: 'producerMetrics',
    initialState,
    reducers: {
        incrementProducerMessage: (state) => {
            state.totalProducerMessages += 1;
            state.intervalProducerMessageCount += 1;
        },
        resetProducerIntervalCounts: (state) => {
            state.intervalProducerMessageCount = 0;
        },
    },
});

export const { incrementProducerMessage, resetProducerIntervalCounts } = producerMetricsSlice.actions;
export default producerMetricsSlice.reducer;
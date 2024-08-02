import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalProducerMessages: 0,
    intervalProducerMessageCount: 0,
};

const producerMetricsSlice = createSlice({
    name: 'producerMetrics',
    initialState,
    reducers: {
        incrementProducerMessage: (state, action) => {
            const count = action.payload || 1;
            state.totalProducerMessages += count;
            state.intervalProducerMessageCount += count;
        },
        resetProducerIntervalCounts: (state) => {
            state.intervalProducerMessageCount = 0;
        },
    },
});

export const { incrementProducerMessage, resetProducerIntervalCounts } = producerMetricsSlice.actions;
export default producerMetricsSlice.reducer;
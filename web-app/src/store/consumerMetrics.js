import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalConsumerReceivedMessages: 0,
    intervalConsumerReceivedMessageCount: 0,
};

const consumerMetricsSlice = createSlice({
    name: 'consumerMetrics',
    initialState,
    reducers: {
        incrementConsumerMessage: (state, action) => {
            const count = action.payload || 1;
            state.totalConsumerReceivedMessages += count;
            state.intervalConsumerReceivedMessageCount += count;
        },
        resetConsumerIntervalCounts: (state) => {
            state.intervalConsumerReceivedMessageCount = 0;
        },
    },
});

export const { incrementConsumerMessage, resetConsumerIntervalCounts } = consumerMetricsSlice.actions;
export default consumerMetricsSlice.reducer;
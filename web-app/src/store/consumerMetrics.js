import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalConsumerReceivedMessages: 0,
    intervalConsumerReceivedMessageCount: 0,
};

const consumerMetricsSlice = createSlice({
    name: 'consumerMetrics',
    initialState,
    reducers: {
        incrementConsumerMessage: (state) => {
            state.totalConsumerReceivedMessages += 1;
            state.intervalConsumerReceivedMessageCount += 1;
        },
        resetConsumerIntervalCounts: (state) => {
            state.intervalConsumerReceivedMessageCount = 0;
        },
    },
});

export const { incrementConsumerMessage, resetConsumerIntervalCounts } = consumerMetricsSlice.actions;
export default consumerMetricsSlice.reducer;
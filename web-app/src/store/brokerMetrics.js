import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalBrokerAcknowledgedMessages: 0,
    intervalBrokerAcknowledgedMessageCount: 0,
};

const brokerMetricsSlice = createSlice({
    name: 'brokerMetrics',
    initialState,
    reducers: {
        addBrokerAcknowledgedMessage: (state) => {
            state.totalBrokerAcknowledgedMessages += 1;
            state.intervalBrokerAcknowledgedMessageCount += 1;
        },
        resetBrokerIntervalCounts: (state) => {
            state.intervalBrokerAcknowledgedMessageCount = 0;
        },
    },
});

export const { addBrokerAcknowledgedMessage, resetBrokerIntervalCounts } = brokerMetricsSlice.actions;
export default brokerMetricsSlice.reducer;
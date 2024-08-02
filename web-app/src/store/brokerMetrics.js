import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    totalBrokerAcknowledgedMessages: 0,
    intervalBrokerAcknowledgedMessageCount: 0,
};

const brokerMetricsSlice = createSlice({
    name: 'brokerMetrics',
    initialState,
    reducers: {
        addBrokerAcknowledgedMessage: (state, action) => {
            const count = action.payload || 1;
            state.totalBrokerAcknowledgedMessages += count;
            state.intervalBrokerAcknowledgedMessageCount += count;
        },
        resetBrokerIntervalCounts: (state) => {
            state.intervalBrokerAcknowledgedMessageCount = 0;
        },
    },
});

export const { addBrokerAcknowledgedMessage, resetBrokerIntervalCounts } = brokerMetricsSlice.actions;
export default brokerMetricsSlice.reducer;
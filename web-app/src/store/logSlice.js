import { createSlice } from '@reduxjs/toolkit';

const initialState = {
    producerLogs: [],
    consumerLogs: [],
    brokerLogs: [],
    proxyLogs: [],
};

const logSlice = createSlice({
    name: 'logs',
    initialState,
    reducers: {
        batchAddLogs: (state, action) => {
            action.payload.forEach(log => {
                if (log.includes('[Producer')) {
                    state.producerLogs.push(log);
                } else if (log.includes('[Consumer')) {
                    state.consumerLogs.push(log);
                } else if (log.includes('[Broker') || log.includes('[Redis')) {
                    state.brokerLogs.push(log);
                } else if (log.includes('[ProxyServer')) {
                    state.proxyLogs.push(log);
                }
            });

            // Trim logs if they exceed a certain length (e.g., keep only the last 100 entries)
            const maxLogEntries = 100;
            state.producerLogs = state.producerLogs.slice(-maxLogEntries);
            state.consumerLogs = state.consumerLogs.slice(-maxLogEntries);
            state.brokerLogs = state.brokerLogs.slice(-maxLogEntries);
            state.proxyLogs = state.proxyLogs.slice(-maxLogEntries);
        },
    },
});

export const { batchAddLogs } = logSlice.actions;
export default logSlice.reducer;
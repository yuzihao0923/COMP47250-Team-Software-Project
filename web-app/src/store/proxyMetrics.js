import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  totalProxyMessages: 0,
  intervalProxyMessageCount: 0,
};

const proxyMetricsSlice = createSlice({
  name: 'proxyMetrics',
  initialState,
  reducers: {
    incrementProxyMessage: (state) => {
      state.totalProxyMessages += 1;
      state.intervalProxyMessageCount += 1;
    },
    resetProxyIntervalCounts: (state) => {
      state.intervalProxyMessageCount = 0;
    },
  },
});

export const { incrementProxyMessage, resetProxyIntervalCounts } = proxyMetricsSlice.actions;
export default proxyMetricsSlice.reducer;
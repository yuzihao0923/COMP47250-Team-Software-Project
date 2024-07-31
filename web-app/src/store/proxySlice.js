import { createSlice } from '@reduxjs/toolkit';

export const proxySlice = createSlice({
  name: 'proxy',
  initialState: {
    proxyLogs: [],
    intervalProxyMessageCount: 0,
    proxyChartData: { labels: [], datasets: [{ data: [] }] }
  },
  reducers: {
    addProxyLog: (state, action) => {
      state.proxyLogs.push(action.payload);
    }
    // resetProxyIntervalCounts: (state) => {
    //   state.intervalProxyMessageCount = 0;
    // },
    // updateProxyChartData: (state, action) => {
    //   state.proxyChartData = action.payload;
    // }
  },
});

export const { addProxyLog } = proxySlice.actions;

export default proxySlice.reducer;

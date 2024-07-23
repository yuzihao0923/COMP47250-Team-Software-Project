import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  warningNumber: 0
};

const settingsSlice = createSlice({
  name: 'settings',
  initialState: initialState,
  reducers: {
    changeWarningNumber(state, action) {
      const { warningNumber } = action.payload;
      state.warningNumber = warningNumber;
    },
    
  },
});

export const { changeWarningNumber } = settingsSlice.actions;
export default settingsSlice.reducer;

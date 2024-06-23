import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  username: '',
  token: '',
  isAuthenticated: false,
};

const userSlice = createSlice({
  name: 'user',
  initialState: initialState,
  reducers: {
    login(state, action) {
      const { user: username, token } = action.payload;
      state.username = username;
      state.token = token;
      state.isAuthenticated = true;
    },
    logout(state) {
      state.username = '';
      state.token = '';
      state.isAuthenticated = false;
    },
  },
});

export const { login, logout } = userSlice.actions;
export default userSlice.reducer;

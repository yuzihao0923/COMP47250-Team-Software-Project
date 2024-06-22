import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    username: '',
    token: ''
}

const userSlice = createSlice({
    name: 'user',
    initialState: initialState,
    reducers: {
        login(state, action) {
            const { user: username, token } = action.payload
            state.username = username
            state.token = token
            // return { ...state, username, token }
        },
        logout(state) {
            // state = initialState
            state.username = ''
            state.token = ''
            // return initialState
        }
    }
})

const { login, logout } = userSlice.actions
const userReducer = userSlice.reducer

export { login, logout }
export default userReducer
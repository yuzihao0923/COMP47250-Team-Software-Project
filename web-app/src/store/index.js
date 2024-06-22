import { configureStore } from '@reduxjs/toolkit'
import userReducer from './userSlice'
import { persistStore, persistReducer } from 'redux-persist'
import storage from 'redux-persist/lib/storage'
import { combineReducers } from "redux"

const persistConfig = {
    key: 'root',   // 在localStorge中生成key为root的值
    storage,
    // blacklist: ['user'],     //设置某个 reducer 数据不持久化，
    whitelist: ['user']        // 设置只有某个 reducer 持久化，其他都不持久化
}

const rootReducer = combineReducers({
    user: userReducer
})

const myPersistReducer = persistReducer(persistConfig, rootReducer)

const store = configureStore({
    reducer: myPersistReducer
})

const persistor = persistStore(store)

export { store, persistor }
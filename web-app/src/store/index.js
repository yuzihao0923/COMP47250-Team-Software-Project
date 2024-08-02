import { configureStore } from '@reduxjs/toolkit';
import userReducer from './userSlice';
import producerMetricsReducer from './producerMetrics';
import consumerMetricsReducer from './consumerMetrics';
import brokerMetricsReducer from './brokerMetrics';
import proxyMetricsReducer from './proxyMetrics';
import logReducer from './logSlice';
import settingsReducer from './settings';
import { persistStore, persistReducer } from 'redux-persist';
import storageSession from 'redux-persist/lib/storage/session';
import { combineReducers } from 'redux';
import { FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER } from 'redux-persist';

const persistConfig = {
  key: 'root',
  storage: storageSession,
  whitelist: ['user', 'producerMetrics', 'consumerMetrics', 'brokerMetrics', 'proxyMetrics', 'settings'],
};

const rootReducer = combineReducers({
  user: userReducer,
  producerMetrics: producerMetricsReducer,
  consumerMetrics: consumerMetricsReducer,
  brokerMetrics: brokerMetricsReducer,
  proxyMetrics: proxyMetricsReducer,
  logs: logReducer,
  settings: settingsReducer
});

const persistedReducer = persistReducer(persistConfig, rootReducer);

const store = configureStore({
  reducer: persistedReducer,
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: [FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER],
      },
    }),
});

const persistor = persistStore(store);

export { store, persistor };
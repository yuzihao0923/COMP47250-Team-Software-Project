import { configureStore } from '@reduxjs/toolkit';
import userReducer from './userSlice';
import producerReducer from './producerSlice';
import consumerReducer from './consumerSlice';
import brokerReducer from './brokerSlice';
import settingsReducer from './settings';
import { persistStore, persistReducer } from 'redux-persist';
import storageSession from 'redux-persist/lib/storage/session';
import { combineReducers } from 'redux';
import { FLUSH, REHYDRATE, PAUSE, PERSIST, PURGE, REGISTER } from 'redux-persist';

const persistConfig = {
  key: 'root',
  storage: storageSession,
  whitelist: ['user', 'producer', 'consumer', 'broker', 'settings'],
};

const rootReducer = combineReducers({
  user: userReducer,
  producer: producerReducer,
  consumer: consumerReducer,
  broker: brokerReducer,
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

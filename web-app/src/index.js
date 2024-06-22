import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import App from './App';
import reportWebVitals from './reportWebVitals';
import { Provider } from 'react-redux'
import { store, persistor } from './store';
import { PersistGate } from 'redux-persist/integration/react';
import { ConfigProvider } from 'antd';
import enUS from 'antd/lib/locale/en_US';
import { ProConfigProvider } from '@ant-design/pro-components';
// import 'antd/dist/antd.css';

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(
  <React.StrictMode>
    <Provider store={store}>
      <PersistGate loading={null} persistor={persistor}>
        <ConfigProvider locale={enUS}>
          <ProConfigProvider value={{ locale: enUS }}>
            <App />
          </ProConfigProvider>
        </ConfigProvider>
      </PersistGate>
    </Provider>
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

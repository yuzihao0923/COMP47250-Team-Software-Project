// src/App.js
import React from 'react';
import { Route, Routes, BrowserRouter } from 'react-router-dom';
import Login from './pages/Login';
import ConsumerConsole from './pages/ConsumerConsole';
import BrokerConsole from './pages/BrokerConsole';
import ProducerConsole from './pages/ProducerConsole';

function App() {
  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Login />} />
          <Route path="/consumer" element={<ConsumerConsole />} />
          <Route path="/broker" element={<BrokerConsole />} />
          <Route path="/producer" element={<ProducerConsole />} />
        </Routes>
      </BrowserRouter>
    </>
  );
}

export default App;

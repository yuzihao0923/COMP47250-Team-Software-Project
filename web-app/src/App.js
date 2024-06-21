import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './pages/Login';
import ConsumerConsole from './pages/ConsumerConsole';
import BrokerConsole from './pages/BrokerConsole';
import ProducerConsole from './pages/ProducerConsole';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/consumer" element={<ConsumerConsole />} />
        <Route path="/broker" element={<BrokerConsole />} />
        <Route path="/producer" element={<ProducerConsole />} />
      </Routes>
    </Router>
  );
}

export default App;

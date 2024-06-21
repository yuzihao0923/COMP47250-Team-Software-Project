import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './components/Login';
import ConsumerConsole from './components/ConsumerConsole';
import BrokerConsole from './components/BrokerConsole';
import ProducerConsole from './components/ProducerConsole';

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

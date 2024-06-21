// src/App.js
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './components/Login';
import BrokerConsole from './components/BrokerConsole';
import './css/App.css';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/broker" element={<BrokerConsole />} />
        <Route exact path="/" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;

// src/App.js
import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './pages/Login.jsx';
import BrokerConsole from './pages/BrokerConsole.jsx';
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
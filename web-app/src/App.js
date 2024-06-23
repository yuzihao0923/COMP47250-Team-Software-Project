import React from 'react';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import Login from './pages/Login';
import BrokerConsole from './pages/BrokerConsole.jsx';
import ProtectedRoute from './components/ProtectRoute.jsx';
import './css/App.css';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/broker" element={<ProtectedRoute><BrokerConsole /></ProtectedRoute>}/>
        <Route exact path="/" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;

import React from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './pages/Login.jsx';
import BrokerConsole from './pages/BrokerConsole.jsx';
import ProtectedRoute from './components/ProtectRoute.jsx';
import Home from './pages/Home.jsx';
import NotFound from './404/NotFound.jsx'
import './css/App.css';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/home" element={<ProtectedRoute><Home /></ProtectedRoute>}>
          {/* <Route index element={<BrokerConsole />} /> */}
          <Route index element={<Navigate to="broker" />} />
          <Route path="broker" element={<BrokerConsole />} />
        </Route>
        <Route exact path="/" element={<Login />} />
        {/* 404 page */}
        <Route path='*' element={<NotFound />} />
      </Routes>
    </Router>
  );
}

export default App;

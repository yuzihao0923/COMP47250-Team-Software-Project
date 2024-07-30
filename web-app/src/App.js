import React from 'react';
import Login from './pages/Login.jsx';
import BrokerConsole from './pages/BrokerConsole.jsx';
import ProtectedRoute from './components/ProtectRoute.jsx';
import Home from './pages/Home.jsx';
import NotFound from './404/NotFound.jsx'
import './css/App.css';
import Register from './pages/Register.jsx';
import Profile from './pages/Profile.jsx';
import Settings from './pages/Settings.jsx';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path='/register' element={<Register />} />
        <Route path="/" element={<ProtectedRoute><Home /></ProtectedRoute>}>
          <Route index element={<BrokerConsole />} />
          <Route path='profile' element={<Profile />} />
          <Route path='settings' element={<Settings />} />
        </Route>
        {/* <Route exact path="/" element={<Login />} /> */}
        {/* 404 page */}
        <Route path='*' element={<NotFound />} />
      </Routes>
    </Router>
  );
}

export default App;

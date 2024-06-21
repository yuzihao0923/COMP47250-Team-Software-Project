// src/components/Login.js
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/api';
import '../css/Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (event) => {
    event.preventDefault();

    try {
      const response = await login({ username, password });
      const { token, role } = response.data;

      if (role === 'broker') {
        localStorage.setItem('token', token);
        navigate('/broker');
      } else {
        setError('this user is not a broker, please try again');
        setUsername('');
        setPassword('');
      }
    } catch (err) {
      if (err.response && err.response.status === 401) {
        setError(err.response.data.message);
        setPassword('');
      } else {
        setError('Login failed. Please try again.');
        setUsername('');
        setPassword('');
      }
    }
  };

  return (
    <div className="login-container">
      <form onSubmit={handleLogin} className="login-form">
        <h2>Broker Login</h2>
        {error && <p className="error">{error}</p>}
        <input
          type="text"
          placeholder="Username"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          required
        />
        <input
          type="password"
          placeholder="Password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;

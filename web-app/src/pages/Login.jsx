import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import '../css/Login.css';

const Login = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async (event) => {
    event.preventDefault();

    try {
      const response = await axios.post('http://localhost:8080/login', {
        username,
        password
      });

      const { token, role } = response.data;

      if (role === 'broker') {
        localStorage.setItem(`${username}_token`, token); // Storing the JWT with username
        navigate('/broker');
      } else {
        setError('this user is not a broker, please try again');
        setUsername('');
        setPassword('');
      }
    } catch (err) {
      if (err.response && err.response.status === 401) {
        const errorMessage = err.response.data;
        setError(errorMessage);
        if (errorMessage.includes('username')) {
          setUsername('');
        }
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

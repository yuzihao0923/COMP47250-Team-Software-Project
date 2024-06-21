import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { login } from '../services/api';
import '../css/Login.css';

const Login = () => {
  const [credentials, setCredentials] = useState({ username: '', password: '', role: 'consumer' });
  const navigate = useNavigate();

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await login(credentials);
      console.log('Login successful:', response.data);

      if (credentials.role === 'consumer') {
        navigate('/consumer');
      } else if (credentials.role === 'broker') {
        navigate('/broker');
      } else if (credentials.role === 'producer') {
        navigate('/producer');
      }
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <div className="login-container">
      <form className="login-form" onSubmit={handleSubmit}>
        <select
          value={credentials.role}
          onChange={(e) => setCredentials({ ...credentials, role: e.target.value })}
        >
          <option value="consumer">Consumer</option>
          <option value="broker">Broker</option>
          <option value="producer">Producer</option>
        </select>
        <input
          type="text"
          value={credentials.username}
          onChange={(e) => setCredentials({ ...credentials, username: e.target.value })}
          placeholder="Username"
        />
        <input
          type="password"
          value={credentials.password}
          onChange={(e) => setCredentials({ ...credentials, password: e.target.value })}
          placeholder="Password"
        />
        <button type="submit">Login</button>
      </form>
    </div>
  );
};

export default Login;

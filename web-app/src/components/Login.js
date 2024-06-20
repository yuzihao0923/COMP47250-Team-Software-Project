import React, { useState } from 'react';
import { login } from '../services/api';

const Login = () => {
  const [credentials, setCredentials] = useState({ username: '', password: '' });

  const handleSubmit = async (event) => {
    event.preventDefault();
    try {
      const response = await login(credentials);
      console.log('Login successful:', response.data);
    } catch (error) {
      console.error('Login failed:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
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
  );
};

export default Login;

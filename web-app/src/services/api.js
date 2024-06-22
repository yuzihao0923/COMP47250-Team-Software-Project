import axios from 'axios';

const API_URL = 'http://localhost:8080';

const getToken = (username) => {
  return localStorage.getItem(`${username}_token`);
};

// Login function
export const login = async (credentials) => {
  try {
    const response = await axios.post(`${API_URL}/login`, credentials);
    const { token, username } = response.data;
    localStorage.setItem(`${username}_token`, token); // store token with username
    return response;
  } catch (error) {
    console.error('Login error:', error);
    throw error;
  }
};

// Function to get logs
export const getLogs = async (username) => {
  try {
    const token = getToken(username);
    const response = await axios.get(`${API_URL}/logs`, {
      headers: {
        'Authorization': `Bearer ${token}`
      }
    });
    return response;
  } catch (error) {
    console.error('Get logs error:', error);
    throw error;
  }
};

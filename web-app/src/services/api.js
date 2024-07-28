import axios from 'axios';

const API_URL = 'http://localhost:8080';

const getToken = (username) => {
  return localStorage.getItem(`${username}_token`);
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

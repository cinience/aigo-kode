import axios from 'axios';

// Export the getConfig function for use in App.jsx
export const getConfig = async () => {
  const response = await axios.get('/api/config');
  return response.data;
};

export const updateConfig = async (config) => {
  const response = await axios.put('/api/config', config);
  return response.data;
};

import axios from 'axios';

const API_BASE_URL = '/api';

// Configuration API
export const getConfig = async () => {
  const response = await axios.get(`${API_BASE_URL}/config`);
  return response.data;
};

export const updateConfig = async (config) => {
  const response = await axios.put(`${API_BASE_URL}/config`, config);
  return response.data;
};

// Chat API
export const sendMessage = async (sessionId, message) => {
  const response = await axios.post(`${API_BASE_URL}/chat`, {
    sessionId,
    message,
  });
  return response.data;
};

export const getChatHistory = async () => {
  const response = await axios.get(`${API_BASE_URL}/chat/history`);
  return response.data;
};

// Tool API
export const executeTool = async (toolName, sessionId, input) => {
  const response = await axios.post(`${API_BASE_URL}/tools/${toolName}`, {
    sessionId,
    input,
  });
  return response.data;
};

// Files API
export const listFiles = async (path = '.') => {
  const response = await axios.get(`${API_BASE_URL}/files`, {
    params: { path },
  });
  return response.data;
};

export const getFileContent = async (path) => {
  const response = await axios.get(`${API_BASE_URL}/files/${encodeURIComponent(path)}`);
  return response.data;
};

export const updateFile = async (path, content) => {
  const response = await axios.put(`${API_BASE_URL}/files/${encodeURIComponent(path)}`, {
    content,
  });
  return response.data;
};

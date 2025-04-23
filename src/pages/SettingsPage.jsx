import React, { useState, useEffect } from 'react';
import { Box, Typography, Paper, TextField, Button, FormControl, InputLabel, Select, MenuItem, Alert } from '@mui/material';
import { getConfig, updateConfig } from '../api';

function SettingsPage() {
  const [config, setConfig] = useState({
    defaultModel: 'gpt-3.5-turbo',
    apiKeys: {
      openai: ''
    }
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await getConfig();
      setConfig({
        defaultModel: data.defaultModel || 'gpt-3.5-turbo',
        apiKeys: {
          openai: '' // We don't receive the actual API key for security reasons
        }
      });
    } catch (err) {
      console.error('Error fetching config:', err);
      setError('Failed to fetch configuration: ' + (err.message || 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleSaveConfig = async () => {
    setLoading(true);
    setError(null);
    setSuccess(false);
    
    try {
      await updateConfig(config);
      setSuccess(true);
    } catch (err) {
      console.error('Error saving config:', err);
      setError('Failed to save configuration: ' + (err.message || 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleModelChange = (e) => {
    setConfig({
      ...config,
      defaultModel: e.target.value
    });
  };

  const handleApiKeyChange = (provider, value) => {
    setConfig({
      ...config,
      apiKeys: {
        ...config.apiKeys,
        [provider]: value
      }
    });
  };

  if (loading && !config) {
    return <Typography>Loading...</Typography>;
  }

  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 3 }}>Settings</Typography>
      
      {error && (
        <Alert severity="error" sx={{ mb: 2 }}>
          {error}
        </Alert>
      )}
      
      {success && (
        <Alert severity="success" sx={{ mb: 2 }}>
          Settings saved successfully!
        </Alert>
      )}
      
      <Paper sx={{ p: 3, mb: 3 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>AI Model Settings</Typography>
        
        <FormControl fullWidth sx={{ mb: 3 }}>
          <InputLabel id="model-select-label">Default Model</InputLabel>
          <Select
            labelId="model-select-label"
            value={config.defaultModel}
            label="Default Model"
            onChange={handleModelChange}
          >
            <MenuItem value="gpt-3.5-turbo">GPT-3.5 Turbo</MenuItem>
            <MenuItem value="gpt-4">GPT-4</MenuItem>
            <MenuItem value="gpt-4-turbo">GPT-4 Turbo</MenuItem>
          </Select>
        </FormControl>
        
        <Typography variant="h6" sx={{ mb: 2 }}>API Keys</Typography>
        
        <TextField
          fullWidth
          label="OpenAI API Key"
          type="password"
          value={config.apiKeys.openai}
          onChange={(e) => handleApiKeyChange('openai', e.target.value)}
          margin="normal"
          helperText="Your API key is stored securely and never shared"
        />
        
        <Box sx={{ mt: 3, display: 'flex', justifyContent: 'flex-end' }}>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleSaveConfig}
            disabled={loading}
          >
            Save Settings
          </Button>
        </Box>
      </Paper>
      
      <Paper sx={{ p: 3 }}>
        <Typography variant="h6" sx={{ mb: 2 }}>About</Typography>
        <Typography variant="body1">
          Go Anon Kode is a terminal-based AI coding tool that can use any model that supports the OpenAI-style API.
          This is a Golang reimplementation of the original anon-kode project.
        </Typography>
        <Typography variant="body2" sx={{ mt: 2, color: 'text.secondary' }}>
          Version: 1.0.0
        </Typography>
      </Paper>
    </Box>
  );
}

export default SettingsPage;

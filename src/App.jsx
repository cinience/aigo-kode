import React, { useState, useEffect } from 'react';
import { Routes, Route } from 'react-router-dom';
import { Box, CssBaseline } from '@mui/material';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import ChatPage from './pages/ChatPage';
import FilesPage from './pages/FilesPage';
import SettingsPage from './pages/SettingsPage';
import { getConfig } from './api/config';

function App() {
  const [drawerOpen, setDrawerOpen] = useState(true);
  const [config, setConfig] = useState({
    defaultModel: '',
    hasCompletedOnboarding: false,
    hasApiKeys: { openai: false }
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchConfig = async () => {
      try {
        const configData = await getConfig();
        setConfig(configData);
      } catch (error) {
        console.error('Failed to fetch config:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchConfig();
  }, []);

  const toggleDrawer = () => {
    setDrawerOpen(!drawerOpen);
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
        Loading...
      </Box>
    );
  }

  // Show settings page if onboarding is not completed or no API keys
  const needsSetup = !config.hasCompletedOnboarding || !config.hasApiKeys.openai;

  return (
    <Box sx={{ display: 'flex' }}>
      <CssBaseline />
      <Header drawerOpen={drawerOpen} toggleDrawer={toggleDrawer} />
      <Sidebar drawerOpen={drawerOpen} />
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          p: 3,
          width: { sm: `calc(100% - ${drawerOpen ? 240 : 0}px)` },
          ml: { sm: `${drawerOpen ? 240 : 0}px` },
          mt: '64px',
        }}
      >
        <Routes>
          <Route path="/" element={needsSetup ? <SettingsPage /> : <ChatPage />} />
          <Route path="/files/*" element={<FilesPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </Box>
    </Box>
  );
}

export default App;

import React, { useState, useEffect, useRef } from 'react';
import { Box, TextField, Button, Typography, Paper, CircularProgress, Divider } from '@mui/material';
import SendIcon from '@mui/icons-material/Send';
import { v4 as uuidv4 } from 'uuid';
import { sendMessage } from '../api';

function ChatPage() {
  const [sessionId] = useState(() => localStorage.getItem('sessionId') || uuidv4());
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const [loading, setLoading] = useState(false);
  const messagesEndRef = useRef(null);

  // Save session ID to localStorage
  useEffect(() => {
    localStorage.setItem('sessionId', sessionId);
  }, [sessionId]);

  // Scroll to bottom when messages change
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSendMessage = async (e) => {
    e.preventDefault();
    if (!input.trim()) return;

    // Add user message to chat
    const userMessage = { role: 'user', content: input };
    setMessages([...messages, userMessage]);
    setInput('');
    setLoading(true);

    try {
      // Send message to API
      const response = await sendMessage(sessionId, input);
      
      // Add assistant message to chat
      const assistantMessage = { role: 'assistant', content: response.response };
      setMessages(prev => [...prev, assistantMessage]);
      
      // Handle tool calls if any
      if (response.toolCalls && response.toolCalls.length > 0) {
        // In a real implementation, we would handle tool calls here
        // For now, we'll just display them in the chat
        const toolCallMessage = {
          role: 'system',
          content: `Tool calls: ${JSON.stringify(response.toolCalls, null, 2)}`,
          isToolCall: true
        };
        setMessages(prev => [...prev, toolCallMessage]);
      }
    } catch (error) {
      console.error('Error sending message:', error);
      const errorMessage = {
        role: 'system',
        content: `Error: ${error.message || 'Failed to send message'}`,
        isError: true
      };
      setMessages(prev => [...prev, errorMessage]);
    } finally {
      setLoading(false);
    }
  };

  const renderMessage = (message, index) => {
    if (message.isToolCall) {
      return (
        <Paper key={index} className="tool-result" elevation={1} sx={{ my: 1, p: 2 }}>
          <Typography variant="body2" component="pre" sx={{ whiteSpace: 'pre-wrap' }}>
            {message.content}
          </Typography>
        </Paper>
      );
    }

    if (message.isError) {
      return (
        <Paper key={index} elevation={1} sx={{ my: 1, p: 2, bgcolor: '#ff5252' }}>
          <Typography variant="body2">{message.content}</Typography>
        </Paper>
      );
    }

    const isUser = message.role === 'user';
    return (
      <Paper 
        key={index} 
        elevation={1} 
        className={isUser ? "message-user" : "message-assistant"}
        sx={{ 
          my: 1, 
          p: 2,
          ml: isUser ? 'auto' : 0,
          mr: isUser ? 0 : 'auto',
          maxWidth: '80%'
        }}
      >
        <Typography variant="body1" sx={{ whiteSpace: 'pre-wrap' }}>
          {message.content}
        </Typography>
      </Paper>
    );
  };

  return (
    <Box sx={{ height: 'calc(100vh - 100px)', display: 'flex', flexDirection: 'column' }}>
      <Typography variant="h5" sx={{ mb: 2 }}>Chat</Typography>
      
      <Box sx={{ flexGrow: 1, overflow: 'auto', mb: 2 }}>
        {messages.length === 0 ? (
          <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100%' }}>
            <Typography variant="body1" color="text.secondary">
              Start a conversation with the AI assistant
            </Typography>
          </Box>
        ) : (
          messages.map(renderMessage)
        )}
        <div ref={messagesEndRef} />
      </Box>
      
      <Divider sx={{ mb: 2 }} />
      
      <Box component="form" onSubmit={handleSendMessage} sx={{ display: 'flex' }}>
        <TextField
          fullWidth
          variant="outlined"
          placeholder="Type your message..."
          value={input}
          onChange={(e) => setInput(e.target.value)}
          disabled={loading}
          sx={{ mr: 1 }}
        />
        <Button 
          type="submit" 
          variant="contained" 
          color="primary" 
          disabled={loading || !input.trim()}
          endIcon={loading ? <CircularProgress size={20} /> : <SendIcon />}
        >
          Send
        </Button>
      </Box>
    </Box>
  );
}

export default ChatPage;

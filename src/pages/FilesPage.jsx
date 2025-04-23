import React, { useState, useEffect } from 'react';
import { Box, Typography, Paper, List, ListItem, ListItemIcon, ListItemText, Breadcrumbs, Link, TextField, Button } from '@mui/material';
import FolderIcon from '@mui/icons-material/Folder';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import ArrowBackIcon from '@mui/icons-material/ArrowBack';
import SaveIcon from '@mui/icons-material/Save';
import { listFiles, getFileContent, updateFile } from '../api';

function FilesPage() {
  const [currentPath, setCurrentPath] = useState('.');
  const [entries, setEntries] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [selectedFile, setSelectedFile] = useState(null);
  const [fileContent, setFileContent] = useState('');
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    fetchFiles(currentPath);
  }, [currentPath]);

  const fetchFiles = async (path) => {
    setLoading(true);
    setError(null);
    try {
      const data = await listFiles(path);
      setEntries(data.entries || []);
    } catch (err) {
      console.error('Error fetching files:', err);
      setError('Failed to fetch files: ' + (err.message || 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleFileClick = async (entry) => {
    if (entry.is_dir) {
      // Navigate to directory
      setCurrentPath(entry.path);
      setSelectedFile(null);
      setFileContent('');
      setIsEditing(false);
    } else {
      // Load file content
      setLoading(true);
      try {
        const data = await getFileContent(entry.path);
        if (data.type === 'text') {
          setSelectedFile(entry);
          setFileContent(data.content || '');
          setIsEditing(false);
        } else if (data.type === 'image') {
          // Handle image files
          setError('Image files are not supported in this view');
        } else {
          setError('Unsupported file type');
        }
      } catch (err) {
        console.error('Error fetching file content:', err);
        setError('Failed to fetch file content: ' + (err.message || 'Unknown error'));
      } finally {
        setLoading(false);
      }
    }
  };

  const handleSaveFile = async () => {
    if (!selectedFile) return;
    
    setLoading(true);
    try {
      await updateFile(selectedFile.path, fileContent);
      setIsEditing(false);
    } catch (err) {
      console.error('Error saving file:', err);
      setError('Failed to save file: ' + (err.message || 'Unknown error'));
    } finally {
      setLoading(false);
    }
  };

  const handleBackClick = () => {
    if (selectedFile) {
      // Go back to file list
      setSelectedFile(null);
      setFileContent('');
      setIsEditing(false);
    } else {
      // Go up one directory
      const parentPath = currentPath.split('/').slice(0, -1).join('/') || '.';
      setCurrentPath(parentPath);
    }
  };

  const renderBreadcrumbs = () => {
    const paths = currentPath === '.' ? [] : currentPath.split('/');
    
    return (
      <Breadcrumbs aria-label="breadcrumb" sx={{ mb: 2 }}>
        <Link 
          color="inherit" 
          href="#" 
          onClick={(e) => {
            e.preventDefault();
            setCurrentPath('.');
          }}
        >
          Root
        </Link>
        {paths.map((path, index) => {
          const pathTo = paths.slice(0, index + 1).join('/');
          return (
            <Link
              key={index}
              color="inherit"
              href="#"
              onClick={(e) => {
                e.preventDefault();
                setCurrentPath(pathTo);
              }}
            >
              {path}
            </Link>
          );
        })}
      </Breadcrumbs>
    );
  };

  return (
    <Box>
      <Typography variant="h5" sx={{ mb: 2 }}>File Explorer</Typography>
      
      <Button 
        startIcon={<ArrowBackIcon />} 
        onClick={handleBackClick}
        disabled={currentPath === '.' && !selectedFile}
        sx={{ mb: 2 }}
      >
        Back
      </Button>
      
      {renderBreadcrumbs()}
      
      {error && (
        <Paper sx={{ p: 2, mb: 2, bgcolor: '#ff5252' }}>
          <Typography>{error}</Typography>
        </Paper>
      )}
      
      {loading ? (
        <Typography>Loading...</Typography>
      ) : selectedFile ? (
        <Box>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
            <Typography variant="h6">{selectedFile.name}</Typography>
            <Box>
              <Button 
                variant={isEditing ? "outlined" : "contained"} 
                color="primary"
                onClick={() => setIsEditing(!isEditing)}
                sx={{ mr: 1 }}
              >
                {isEditing ? "Cancel" : "Edit"}
              </Button>
              {isEditing && (
                <Button 
                  variant="contained" 
                  color="primary"
                  startIcon={<SaveIcon />}
                  onClick={handleSaveFile}
                >
                  Save
                </Button>
              )}
            </Box>
          </Box>
          
          {isEditing ? (
            <TextField
              fullWidth
              multiline
              rows={20}
              value={fileContent}
              onChange={(e) => setFileContent(e.target.value)}
              variant="outlined"
            />
          ) : (
            <Paper className="code-block" sx={{ p: 2 }}>
              <pre style={{ margin: 0, whiteSpace: 'pre-wrap', wordBreak: 'break-word' }}>
                {fileContent}
              </pre>
            </Paper>
          )}
        </Box>
      ) : (
        <List component={Paper}>
          {entries.length === 0 ? (
            <ListItem>
              <ListItemText primary="No files found" />
            </ListItem>
          ) : (
            entries.map((entry) => (
              <ListItem 
                button 
                key={entry.path} 
                onClick={() => handleFileClick(entry)}
              >
                <ListItemIcon>
                  {entry.is_dir ? <FolderIcon /> : <InsertDriveFileIcon />}
                </ListItemIcon>
                <ListItemText 
                  primary={entry.name} 
                  secondary={entry.is_dir ? 'Directory' : `${entry.size} bytes`} 
                />
              </ListItem>
            ))
          )}
        </List>
      )}
    </Box>
  );
}

export default FilesPage;

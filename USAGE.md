#!/bin/bash

# Usage guide for Go Anon Kode
# This script provides examples of how to use the Go Anon Kode CLI and server

echo "Go Anon Kode Usage Guide"
echo "========================"
echo ""

echo "CLI Usage Examples:"
echo "------------------"
echo "# Start interactive CLI mode"
echo "./go-anon-kode-cli"
echo ""
echo "# Read a file"
echo "./go-anon-kode-cli file-read /path/to/file.txt"
echo ""
echo "# Write to a file"
echo "./go-anon-kode-cli file-write /path/to/file.txt \"Content to write\""
echo ""
echo "# Execute a bash command"
echo "./go-anon-kode-cli bash \"ls -la\""
echo ""
echo "# Search for files"
echo "./go-anon-kode-cli glob \"*.go\" /path/to/search"
echo ""
echo "# Search file content"
echo "./go-anon-kode-cli grep \"search pattern\" /path/to/file.txt"
echo ""

echo "Server Usage:"
echo "------------"
echo "# Start the server"
echo "./go-anon-kode-server"
echo ""
echo "# Access the web interface"
echo "Open http://localhost:8080 in your browser"
echo ""

echo "Configuration:"
echo "-------------"
echo "# API keys and settings are stored in:"
echo "~/.go-anon-kode/config.json"
echo ""
echo "# You can configure the API key via the web interface"
echo "# or by editing the config file directly"
echo ""

echo "For more information, see the README.md file."

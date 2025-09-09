#!/bin/bash

# Axon Visual Editor Startup Script
echo "🚀 Starting Axon Visual Editor..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.24 or later."
    exit 1
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18 or later."
    exit 1
fi

# Check if gopls is available (optional but recommended)
if ! command -v gopls &> /dev/null; then
    echo "⚠️  gopls not found. Installing for better Go autocompletion..."
    go install golang.org/x/tools/gopls@latest
fi

# Build frontend if needed
if [ ! -d "frontend/dist" ]; then
    echo "📦 Building frontend..."
    cd frontend
    npm install
    npm run build
    cd ..
fi

# Start the backend server
echo "🔧 Starting backend server..."
cd backend
go run *.go &
BACKEND_PID=$!
cd ..

# Wait a moment for the server to start
sleep 3

echo ""
echo "✅ Axon Visual Editor is running!"
echo ""
echo "🌐 Open your browser and go to: http://localhost:8080"
echo "📚 API documentation available at: http://localhost:8080/api/health"
echo ""
echo "Features available:"
echo "  ✓ Visual node-based graph editor"
echo "  ✓ Real-time graph validation"
echo "  ✓ Go code transpilation"
echo "  ✓ Graph loading/saving"
echo "  ✓ Go language autocompletion (gopls integration)"
echo ""
echo "Press Ctrl+C to stop the server..."

# Function to cleanup on exit
cleanup() {
    echo ""
    echo "🛑 Stopping Axon Visual Editor..."
    kill $BACKEND_PID 2>/dev/null
    exit 0
}

# Set trap to cleanup on Ctrl+C
trap cleanup INT

# Wait for the backend process
wait $BACKEND_PID
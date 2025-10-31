#!/bin/bash
# Installation script for Backend.im CLI

set -e

echo "🔨 Building Backend.im CLI binary..."

# Build binary using Docker (no local Go installation needed)
cd "$(dirname "$0")/cli"
docker run --rm -v "$(pwd)":/app -w /app golang:1.21-alpine sh -c \
  "go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend-im ./cmd/backend-im"

# Fix permissions
sudo chown $USER:$USER backend-im 2>/dev/null || chown $USER:$USER backend-im
chmod +x backend-im

# Install to /usr/local/bin (or suggest user path)
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
BINARY_PATH="$INSTALL_DIR/backend-im"

echo ""
echo "📦 Installing to $INSTALL_DIR..."

if [ -w "$INSTALL_DIR" ]; then
    cp backend-im "$BINARY_PATH"
    echo "✅ Installed successfully to $BINARY_PATH"
else
    echo "⚠️  Need sudo to install to $INSTALL_DIR"
    sudo cp backend-im "$BINARY_PATH"
    echo "✅ Installed successfully to $BINARY_PATH"
fi

# Create config directory (with proper permissions)
if [ ! -d ~/.backend-im ]; then
    mkdir -p ~/.backend-im
    chmod 700 ~/.backend-im 2>/dev/null || sudo chown $USER:$USER ~/.backend-im && chmod 700 ~/.backend-im
fi

echo ""
echo "✅ Installation complete!"
echo ""
echo "💡 To use with local mock API, set:"
echo "   export BACKEND_IM_API_URL=http://localhost:8080"
echo ""
echo "   Or add to your ~/.bashrc or ~/.zshrc:"
echo "   echo 'export BACKEND_IM_API_URL=http://localhost:8080' >> ~/.bashrc"
echo ""
echo "🚀 Usage:"
echo "   backend-im login"
echo "   backend-im generate \"Your prompt\" --project my-api"
echo "   backend-im commit my-api"
echo "   backend-im deploy my-api"


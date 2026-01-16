#!/bin/bash

# Swagger Setup and Generation Script

echo "🔧 Setting up Swagger for Construction Backend..."

# Install swag CLI if not already installed
if ! command -v swag &> /dev/null; then
    echo "📦 Installing swag CLI..."
    go install github.com/swaggo/swag/cmd/swag@latest
    
    # Add Go bin to PATH if not already there
    export PATH=$PATH:$(go env GOPATH)/bin
fi

# Generate Swagger documentation
echo "📝 Generating Swagger documentation..."
swag init -g cmd/api/main.go -o ./docs

if [ $? -eq 0 ]; then
    echo "✅ Swagger documentation generated successfully!"
    echo ""
    echo "🚀 Access Swagger UI at:"
    echo "   http://localhost:8080/swagger/index.html"
    echo ""
    echo "📖 Your API documentation is ready!"
else
    echo "❌ Failed to generate Swagger documentation"
    exit 1
fi

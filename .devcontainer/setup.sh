#!/bin/bash
set -e

echo "🚀 Setting up Search API development environment..."

# Install Dagger
echo "📦 Installing Dagger CLI..."
cd /usr/local
curl -L https://dl.dagger.io/dagger/install.sh | DAGGER_VERSION=0.9.5 sh
cd -

# Verify Dagger installation
dagger version

# Restore .NET dependencies
echo "📦 Restoring .NET dependencies..."
dotnet restore SearchApi.sln

# Install dotnet tools
echo "🔧 Installing .NET tools..."
dotnet tool restore 2>/dev/null || dotnet new tool-manifest && dotnet tool install dotnet-format

# Initialize Dagger module
echo "🔧 Initializing Dagger module..."
cd .dagger
go mod download
cd -

# Make scripts executable
chmod +x scripts/*.sh

# Create local environment file
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
fi

echo ""
echo "✅ Setup complete!"
echo ""
echo "Available commands:"
echo "  make dev          - Start development environment with Docker Compose"
echo "  make test         - Run unit tests"
echo "  make dagger-full  - Run full Dagger CI/CD pipeline"
echo "  dagger call --help - See all available Dagger functions"
echo ""
echo "To start developing:"
echo "  1. Run 'make dev' to start Solr"
echo "  2. Open SearchApi/Program.cs"
echo "  3. Press F5 to debug, or run 'dotnet watch run' in SearchApi/"
echo ""

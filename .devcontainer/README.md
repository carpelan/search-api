# GitHub Codespaces Quick Start

This repository is fully configured for GitHub Codespaces with all dependencies pre-installed.

## 🚀 Launch Codespace

1. Click the **Code** button on GitHub
2. Select **Codespaces** tab
3. Click **Create codespace on claude/dagger-csharp-ci-pipeline-...**

The environment will automatically set up with:
- .NET 8.0 SDK
- Dagger CLI
- Docker-in-Docker
- Go 1.21
- kubectl & Helm

## 📦 What's Included

The codespace includes:
- ✅ .NET 8.0 development environment
- ✅ Dagger for CI/CD pipeline execution
- ✅ Docker for container builds
- ✅ All necessary VS Code extensions
- ✅ Port forwarding for API (8080) and Solr (8983)

## 🎯 Quick Commands

### Development
\`\`\`bash
# Start Solr and development environment
make dev

# Run the API
cd SearchApi && dotnet watch run

# Run tests
make test
\`\`\`

### Dagger Pipeline
\`\`\`bash
# Build and test
dagger call build --source=.

# Run security scans
dagger call security-scan --source=.
dagger call generate-sbom --source=.

# Build container
dagger call build-container --source=.

# Full pipeline
dagger call full-pipeline --source=. --tag=latest
\`\`\`

### Local Development
\`\`\`bash
# Use Dagger for everything
dagger call full-pipeline --source=.

# Or start just Solr for local .NET development
docker run -d -p 8983:8983 solr:9.4 solr-precreate metadata
cd SearchApi && dotnet run
\`\`\`

## 🔍 Accessing Services

Once the codespace is running:

- **Search API**: Will be available at the forwarded port 8080
- **Solr Admin**: Available at port 8983 (http://localhost:8983/solr)
- **Swagger UI**: http://localhost:8080/swagger

## 📝 Environment Configuration

The `.env` file is automatically created from `.env.example`. Update it with:

\`\`\`bash
# Solr Configuration
SOLR_URL=http://solr:8983/solr/metadata

# For Harbor registry (optional)
HARBOR_URL=your-harbor-url
HARBOR_USERNAME=your-username
HARBOR_PASSWORD=your-password
\`\`\`

## 🧪 Testing OAI-PMH Harvesting

Once the environment is running:

\`\`\`bash
# Start the API and Solr
make dev

# In another terminal, trigger a harvest
curl -X POST "http://localhost:8080/api/harvest/oai-pmh?metadataPrefix=oai_dc&maxRecords=10"

# Search the indexed data
curl -X POST "http://localhost:8080/api/search/search" \
  -H "Content-Type: application/json" \
  -d '{"query": "*:*", "rows": 10}'
\`\`\`

## 🐛 Debugging

The codespace includes the C# extension, so you can:

1. Open `SearchApi/Program.cs`
2. Set breakpoints
3. Press F5 to start debugging

## 📚 Documentation

- [Dagger Documentation](https://docs.dagger.io)
- [.NET 8 Documentation](https://learn.microsoft.com/en-us/dotnet/)
- [Apache Solr Documentation](https://solr.apache.org/guide/)
- [Riksarkivet APIs](https://sok.riksarkivet.se/data-api)

## 💡 Tips

- The codespace includes Docker-in-Docker, so all Docker commands work
- Dagger uses Docker under the hood for reproducible builds
- Ports 8080, 8983 are automatically forwarded
- The terminal supports multiple tabs for running services concurrently

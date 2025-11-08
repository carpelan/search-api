# Dagger CI/CD Pipeline Demonstration

**A comprehensive showcase of modern CI/CD practices using Dagger**

This repository demonstrates how to build a complete, security-first CI/CD pipeline using [Dagger](https://dagger.io) with a real-world C# application. The example application is a Search API for Riksarkivet metadata, but the **primary focus is demonstrating Dagger's capabilities** for building reproducible, containerized CI/CD workflows.

> **Note**: This is a proof-of-concept demonstrating CI/CD patterns with Dagger. The Search API serves as a realistic example application to showcase the pipeline, not as a production search solution.

## 🌟 What This Demonstrates

### Dagger CI/CD Pipeline Features (Primary Focus)

This project showcases a **12-step CI/CD pipeline** built entirely with Dagger:
1. ✅ **Automated Build & Test** - .NET 8.0 compilation and unit tests
2. ✅ **Static Code Analysis** - Code formatting and quality checks with dotnet format
3. ✅ **Security Dependency Scanning** - Vulnerable package detection
4. ✅ **SBOM Generation** - Software Bill of Materials using Syft
5. ✅ **Multi-Stage Container Build** - Optimized Docker images
6. ✅ **Container Security Scanning** - Trivy vulnerability analysis
7. ✅ **Local Registry Push** - Testing registry workflow
8. ✅ **K3s Cluster Provisioning** - Ephemeral Kubernetes cluster in Dagger
9. ✅ **Service Deployment** - Automated Solr deployment to K3s
10. ✅ **Application Deployment** - API deployment with health checks
11. ✅ **Integration Testing** - End-to-end tests in live cluster
12. ✅ **Production Registry Push** - Multi-registry support (Harbor, GHCR, Docker Hub)

### Example Application (Search API)

The pipeline demonstrates these practices on a real C# application:
- **RESTful API** with Swagger/OpenAPI
- **Solr Integration** for full-text search
- **OAI-PMH Integration** for metadata harvesting (Riksarkivet)
- **Security-hardened containers** (non-root, minimal attack surface)
- **Structured logging** with Serilog

## 🏗️ Architecture

```
┌─────────────────┐
│   Search API    │
│   (.NET 8.0)    │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Solr 9.4      │
│  (Full-text     │
│   search)       │
└─────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- [Dagger](https://docs.dagger.io/install) installed
- Docker or Podman
- .NET 8.0 SDK (for local development)

### Run the Full Pipeline

```bash
# Run the complete CI/CD pipeline
dagger call full-pipeline --source=.

# With Harbor registry push
dagger call full-pipeline \
  --source=. \
  --harbor-url=harbor.example.com \
  --harbor-username=env:HARBOR_USERNAME \
  --harbor-password=env:HARBOR_PASSWORD \
  --harbor-project=search-api \
  --tag=v1.0.0
```

### Individual Pipeline Steps

```bash
# Build and run unit tests
dagger call build --source=.

# Run security scan
dagger call security-scan --source=.

# Run static analysis
dagger call static-analysis --source=.

# Generate SBOM
dagger call generate-sbom --source=.

# Build container
dagger call build-container --source=.

# Scan container for vulnerabilities
dagger call scan-container \
  --container=$(dagger call build-container --source=.)

# Setup K3s cluster for testing
dagger call setup-k3s

# Run integration tests
dagger call run-integration-tests \
  --source=. \
  --cluster=$(dagger call setup-k3s)
```

## 📋 CI/CD Pipeline Steps

The full pipeline executes the following steps:

1. **Build & Unit Test** - Compile C# code and run unit tests
2. **Static Analysis** - Code formatting and quality verification
3. **Security Dependency Scan** - Check for vulnerable packages
4. **SBOM Generation** - Create Software Bill of Materials
5. **Container Build** - Multi-stage Docker build
6. **Container Security Scan** - Trivy vulnerability scanning
7. **Local Registry Push** - Push to test registry
8. **K3s Cluster Setup** - Provision local Kubernetes cluster
9. **Solr Deployment** - Deploy and configure Solr
10. **API Deployment** - Deploy Search API to K3s
11. **Integration Tests** - End-to-end testing
12. **Harbor Registry Push** - Push production image (optional)

## 🔒 Security Features

### Shift-Left Security Practices

1. **Dependency Scanning** - Automated vulnerability detection in NuGet packages
2. **Static Analysis** - Code quality and security pattern detection
3. **SBOM Generation** - Complete software bill of materials
4. **Container Scanning** - Multi-layer container vulnerability analysis
5. **Non-root Execution** - Containers run as unprivileged users
6. **Secret Management** - Integration with Infisical
7. **Network Policies** - Kubernetes network segmentation
8. **Resource Limits** - CPU and memory constraints

### Container Security

- Base images: Official Microsoft .NET images
- Multi-stage builds to minimize attack surface
- Non-root user execution
- No unnecessary packages
- Security scanning with Trivy
- SBOM generation with Syft

## 🎯 API Endpoints

### Search Operations

```bash
# Search documents
POST /api/search/search
{
  "query": "riksarkivet",
  "rows": 10,
  "start": 0,
  "sortField": "created_date",
  "sortOrder": "desc"
}

# Get document by ID
GET /api/search/{id}

# Index new document
POST /api/search/index
{
  "id": "doc-123",
  "title": "Document Title",
  "description": "Document description",
  ...
}

# Delete document
DELETE /api/search/{id}

# Health check
GET /health
```

### Harvest Riksarkivet Data via OAI-PMH

This API is designed to work with metadata harvested from Riksarkivet's OAI-PMH service.

**OAI-PMH Endpoint**: `https://oai-pmh.riksarkivet.se/OAI`

For harvesting and integration details, see:
- 📖 **[OAI-PMH Integration Guide](docs/OAI-PMH-INTEGRATION.md)**
- 🔧 **[Riksarkivet Dataplattform](https://github.com/Riksarkivet/dataplattform)**
- 📚 **[Official OAI-PMH Wiki](https://github.com/Riksarkivet/dataplattform/wiki/OAI-PMH)**

Quick example - List available collections:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListAllAuth"
```

Harvest records:
```bash
curl "https://oai-pmh.riksarkivet.se/OAI?verb=ListRecords&metadataPrefix=oai_ra_ead" > records.xml
```

Then parse the EAD XML and POST to this API's `/api/search/index` endpoint.

## 🧪 Testing

### Unit Tests

```bash
dotnet test SearchApi.Tests/SearchApi.Tests.csproj
```

### Integration Tests

```bash
# Tests run against live Solr instance
SOLR_URL=http://solr:8983/solr/metadata \
dotnet test SearchApi.IntegrationTests/SearchApi.IntegrationTests.csproj
```

## 🐳 Docker Deployment

### Build Locally

```bash
docker build -t search-api:latest .
```

### Run with Docker Compose

```bash
version: '3.8'
services:
  solr:
    image: solr:9.4
    ports:
      - "8983:8983"
    command:
      - solr-precreate
      - metadata

  search-api:
    image: search-api:latest
    ports:
      - "8080:8080"
    environment:
      - Solr__Url=http://solr:8983/solr/metadata
    depends_on:
      - solr
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes

```bash
# Create namespace and deploy Solr
kubectl apply -f k8s/solr-deployment.yaml

# Deploy Search API
kubectl apply -f k8s/api-deployment.yaml

# Port forward to access locally
kubectl port-forward -n search-system svc/search-api 8080:80
```

### Access the API

```bash
# Check health
curl http://localhost:8080/health

# Search
curl -X POST http://localhost:8080/api/search/search \
  -H "Content-Type: application/json" \
  -d '{"query": "*:*", "rows": 10}'
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ASPNETCORE_URLS` | URLs to bind to | `http://+:8080` |
| `Solr__Url` | Solr connection URL | `http://solr:8983/solr/metadata` |
| `ASPNETCORE_ENVIRONMENT` | Environment name | `Production` |

### Solr Configuration

The Solr schema includes the following fields:

- `id` - Unique document identifier
- `title` - Document title (text, indexed)
- `description` - Document description (text, indexed)
- `author` - Author name (string, indexed)
- `created_date` - Creation date (date, indexed)
- `modified_date` - Modification date (date, indexed)
- `tags` - Tags (multi-valued strings)
- `content_type` - MIME type
- `file_size` - File size in bytes
- `full_text` - Full document text (text, indexed)

## 📊 Monitoring

### Metrics

- Health checks at `/health`
- Solr admin UI at `http://solr:8983/solr/#/`
- API documentation at `/swagger`

### Logs

Structured JSON logs with:
- Request/response details
- Search query performance
- Error tracking
- Security events

## 🔄 CI/CD Integration

### Argo Workflows

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: search-api-ci-
spec:
  entrypoint: build-and-deploy
  templates:
  - name: build-and-deploy
    steps:
    - - name: dagger-pipeline
        template: dagger
  - name: dagger
    container:
      image: dagger:latest
      command: [dagger]
      args:
        - call
        - full-pipeline
        - --source=/workspace
        - --harbor-url={{workflow.parameters.harbor-url}}
        - --tag={{workflow.parameters.tag}}
```

### Argo CD

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: search-api
spec:
  project: default
  source:
    repoURL: https://github.com/your-org/search-api
    path: k8s
    targetRevision: main
  destination:
    server: https://kubernetes.default.svc
    namespace: search-system
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

## 📚 Data Source

This API is designed to index and search metadata harvested from **Riksarkivet's OAI-PMH service**.

### Riksarkivet OAI-PMH

- **Endpoint**: `https://oai-pmh.riksarkivet.se/OAI`
- **Documentation**: https://github.com/Riksarkivet/dataplattform/wiki/OAI-PMH
- **Metadata Formats**: `oai_ape_ead` (Archives Portal Europe), `oai_ra_ead` (Riksarkivet EAD)
- **License**: CC0 1.0 Universal

### Integration

See [OAI-PMH Integration Guide](docs/OAI-PMH-INTEGRATION.md) for detailed instructions on:
- Harvesting metadata from Riksarkivet
- Parsing EAD XML format
- Indexing into this search API
- Using existing Riksarkivet tools

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run the full pipeline: `dagger call full-pipeline --source=.`
5. Submit a pull request

## 📄 License

MIT License - See LICENSE file for details

## 🆘 Support

- Documentation: See `/swagger` endpoint
- Issues: GitHub Issues
- Security: See SECURITY.md for reporting vulnerabilities

## 🎉 Acknowledgments

- Inspired by [AI-Riksarkivet/coder-templates](https://github.com/AI-Riksarkivet/coder-templates)
- Built with [Dagger](https://dagger.io)
- Powered by [.NET 8.0](https://dotnet.microsoft.com/) and [Apache Solr](https://solr.apache.org/)

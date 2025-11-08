# Dagger CI/CD Pipeline - Security-First Demonstration

**Showcasing shift-left security practices with Dagger**

This repository demonstrates how to build a **security-first CI/CD pipeline** using [Dagger](https://dagger.io). The focus is on integrating comprehensive security scanning and enforcement into every stage of the pipeline using Dagger's programmable approach.

The example application is a C# Search API, chosen to demonstrate security practices across:
- Dependency vulnerability scanning
- Container security
- Infrastructure-as-Code validation
- SBOM generation
- Policy enforcement

> **Primary Focus**: Security automation in CI/CD using Dagger
>
> **Secondary**: The Search API is a realistic example to showcase security scanning on a multi-component application (API + Solr + Kubernetes)

## 🔒 Security-First Pipeline (12 Steps)

This demonstrates a **security-focused CI/CD pipeline** with Dagger:

### Security Stages
1. ✅ **Build & Unit Test** - Compilation with security in mind
2. ⚠️ **Static Code Analysis** - Code formatting (NOT full SAST - see gaps below)
3. ⚠️ **Dependency Vulnerability Scan** - NuGet packages (no enforcement - see gaps)
4. ✅ **SBOM Generation** - Complete software bill of materials (Syft)
5. ✅ **Container Build** - Multi-stage, non-root user
6. ⚠️ **Container Security Scan** - Trivy (currently doesn't fail build - see gaps)
7. ✅ **Local Registry Push** - Secure image distribution
8. ✅ **K3s Cluster** - Ephemeral test environment
9. ✅ **Service Deployment** - Solr with security context
10. ✅ **API Deployment** - Non-root, resource-limited containers
11. ✅ **Integration Testing** - Security validation
12. ✅ **Registry Push** - Multi-registry support (Harbor, GHCR)

### 🚨 Known Security Gaps (Intentional for demonstration)

This POC demonstrates basic security scanning but has **intentional gaps** to show what a production pipeline needs:

- ❌ **No secret scanning** (GitLeaks, TruffleHog)
- ❌ **No real SAST** (current static analysis is just code formatting)
- ❌ **No policy enforcement** (scans run but don't fail builds)
- ❌ **No IaC scanning** (Kubernetes manifests not checked)
- ❌ **No image signing** (Cosign, Sigstore)
- ❌ **No license compliance**

**See [SECURITY-CI-ANALYSIS.md](docs/SECURITY-CI-ANALYSIS.md)** for comprehensive analysis and recommended improvements.

### Why This Application?

The C# Search API was chosen to demonstrate security across multiple attack surfaces:
- **Application code** - .NET vulnerabilities (XSS, injection, etc.)
- **Dependencies** - Third-party package vulnerabilities
- **Container images** - OS and runtime vulnerabilities
- **Infrastructure** - Kubernetes misconfigurations
- **External integrations** - OAI-PMH, Solr connections
- **Secrets management** - Database credentials, API keys

This complexity showcases where security scanning fits in a realistic CI/CD pipeline.

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

- [Dagger](https://docs.dagger.io/install) installed (includes Docker/Podman engine)
- .NET 8.0 SDK (optional, for local development outside Dagger)

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

## 📋 Understanding This Demo

**Important**: This repository is a **Dagger CI/CD demonstration**, not primarily a search API project.

### What This Is

✅ **A showcase of Dagger CI/CD capabilities**
- How to build portable, reproducible pipelines
- Security-first development practices
- Local development parity with CI
- Multi-platform compatibility (works on any CI system)

✅ **A realistic example application**
- C# .NET 8.0 API demonstrating real-world complexity
- Integration with external services (Solr, OAI-PMH)
- Kubernetes deployment patterns
- Security hardening practices

### What This Is NOT

❌ A production-ready search solution
❌ A comprehensive Solr tutorial
❌ A Riksarkivet integration guide (see their official tools)
❌ A Kubernetes best practices reference

### For More Details

See **[DAGGER-SHOWCASE.md](docs/DAGGER-SHOWCASE.md)** for:
- Why we chose Dagger
- Key patterns demonstrated
- How to adapt this for your projects
- Integration with various CI platforms

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

## 🚀 Local Development with Dagger

All containerization is handled by Dagger. No need for separate Dockerfiles or docker-compose!

### Build and Run with Dagger

```bash
# Build the application
dagger call build --source=.

# Build container image
dagger call build-container --source=.

# Run full pipeline locally (includes Solr + API in K3s)
dagger call full-pipeline --source=.
```

### Run API Locally (Outside Dagger)

If you want to run the .NET API directly for development:

```bash
# Start Solr using Dagger
dagger call deploy-solr --k8s-manifests=k8s --cluster=$(dagger call setup-k3s)

# Or use a simple Solr container
docker run -p 8983:8983 solr:9.4 solr-precreate metadata

# Run the API
cd SearchApi && dotnet run
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

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

## 🔒 Comprehensive Security-First Pipeline (24 Steps)

This demonstrates a **production-grade security-focused CI/CD pipeline** with Dagger implementing **9 enforced security gates**:

### 🛡️ Security Gates (Fail-Fast)

**GATE 1: 🔐 Secret Scanning** - TruffleHog detects hardcoded secrets (BLOCKS pipeline)
**GATE 2: 🛡️ SAST (Generic)** - Semgrep finds security vulnerabilities in code (BLOCKS pipeline)
**GATE 3: 🔒 SAST (C# Specific)** - .NET Analyzers for C#-specific security issues (BLOCKS pipeline)
**GATE 4: 🔒 Dependency Scan** - Trivy checks for vulnerable packages (BLOCKS pipeline)
**GATE 5: 📜 License Compliance** - Trivy detects problematic licenses (BLOCKS pipeline)
**GATE 6: 📐 Policy as Code** - OPA/Conftest validates configurations against custom policies (BLOCKS pipeline)
**GATE 7: 🔎 Container Scan** - Trivy blocks HIGH/CRITICAL vulnerabilities (BLOCKS pipeline)
**GATE 8: 🎯 DAST** - OWASP ZAP tests running application for vulnerabilities (BLOCKS pipeline)
**GATE 9: 🔓 API Security** - Nuclei tests for OWASP API Top 10 vulnerabilities (BLOCKS pipeline)

### Complete Pipeline Steps

1. ✅ **Secret Scanning** - TruffleHog (enforced, fails on secrets)
2. ✅ **SAST (Generic)** - Semgrep security analysis (enforced, fails on vulnerabilities)
3. ✅ **SAST (C# Specific)** - .NET Security Analyzers (enforced, 400+ rules)
4. ✅ **Build & Unit Test** - Compilation and testing
5. ✅ **Code Coverage** - XPlat Coverage with 80% threshold
6. ✅ **Code Quality** - dotnet format validation
7. ✅ **Dependency Scan** - Trivy filesystem scan (enforced, fails on HIGH/CRITICAL)
8. ✅ **License Compliance** - Trivy license scan (enforced, blocks problematic licenses)
9. ✅ **IaC Security** - Checkov for Kubernetes manifests
10. ✅ **Policy as Code** - OPA/Conftest validates K8s configurations (enforced, fails on policy violations)
11. ✅ **SBOM Generation** - Syft generates software bill of materials
12. ✅ **Container Build** - Multi-stage, non-root user
13. ✅ **Container Scan** - Trivy image scan (enforced, fails on HIGH/CRITICAL)
14. ✅ **CIS Benchmark** - Docker CIS compliance validation (enforced, reports HIGH/CRITICAL)
15. ✅ **SBOM Attestation** - Cosign attaches signed SBOM to image
16. ✅ **Registry Push** - Local registry for testing
17. ✅ **K3s Cluster** - Ephemeral test environment
18. ✅ **Solr Deployment** - Database with security context
19. ✅ **API Deployment** - Non-root, resource-limited containers
20. ✅ **Integration Tests** - End-to-end validation
21. ✅ **DAST** - OWASP ZAP dynamic security testing (enforced, fails on vulnerabilities)
22. ✅ **API Security Testing** - Nuclei scans for OWASP API Top 10 (enforced)
23. ✅ **Performance Testing** - k6 load tests (optional, configurable thresholds)
24. ✅ **Mutation Testing** - Stryker.NET tests test quality (optional, can be slow)

### 🎯 Security Features Implemented

**Shift-Left Security** ✅
- ✅ Secret scanning with enforcement (TruffleHog)
- ✅ SAST with enforcement (Semgrep + .NET Analyzers) - static code analysis
- ✅ DAST with enforcement (OWASP ZAP) - dynamic runtime testing
- ✅ API Security testing (Nuclei) - OWASP API Top 10
- ✅ Dependency vulnerability scanning with enforcement (Trivy)
- ✅ License compliance scanning with enforcement (Trivy)
- ✅ Container vulnerability scanning with enforcement (Trivy)
- ✅ IaC security scanning (Checkov)
- ✅ Policy as Code enforcement (OPA/Conftest)
- ✅ CIS Benchmark compliance validation (Trivy)
- ✅ SBOM generation (Syft)
- ✅ SBOM attestation with cryptographic signing (Cosign)
- ✅ Non-root container execution
- ✅ Resource limits and security contexts

**Supply Chain Security** ✅
- ✅ Complete dependency tracking
- ✅ Multi-layer vulnerability detection
- ✅ SBOM in SPDX format
- ✅ SBOM attestation with cryptographic signing
- ✅ License compliance enforcement
- ✅ Image signing capability (Cosign/Sigstore)
- ✅ Secure container registry integration
- ✅ CIS Docker Benchmark compliance

**Runtime Security** ✅
- ✅ Dynamic security testing against live application
- ✅ OWASP Top 10 vulnerability detection
- ✅ OWASP API Security Top 10 testing
- ✅ XSS, SQLi, auth bypass detection
- ✅ API security testing (Nuclei)

**Quality & Performance** ✅
- ✅ Code coverage enforcement (80% threshold)
- ✅ Performance testing with k6
- ✅ Mutation testing with Stryker.NET
- ✅ Load testing with configurable thresholds

### 📊 Security Enforcement Policy

| Check Type | Tool | Severity Threshold | Action |
|------------|------|-------------------|--------|
| Secrets | TruffleHog | Any | **FAIL** |
| Code Vulnerabilities (SAST Generic) | Semgrep | ERROR, WARNING | **FAIL** |
| Code Vulnerabilities (SAST C#) | .NET Analyzers | Any | **FAIL** |
| Dependencies | Trivy | HIGH, CRITICAL | **FAIL** |
| License Compliance | Trivy | HIGH, CRITICAL | **FAIL** |
| Policy Violations | OPA/Conftest | Any | **FAIL** |
| Container | Trivy | HIGH, CRITICAL | **FAIL** |
| CIS Benchmark | Trivy | HIGH, CRITICAL | Report |
| Runtime Vulnerabilities (DAST) | OWASP ZAP | Any | **FAIL** |
| API Security | Nuclei | HIGH, CRITICAL | **FAIL** |
| Code Coverage | XPlat Coverage | <80% | **FAIL** |
| IaC | Checkov | INFO | Report |
| Performance | k6 | Configurable | Warn |
| Mutation Score | Stryker.NET | <80% | Warn |

**Result**: Vulnerable code cannot reach production - tested statically, dynamically, AND for API-specific vulnerabilities.

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
# Run the complete CI/CD pipeline (source defaults to current directory)
dagger call full-pipeline

# With container registry push (works with any registry)
# Example 1: Harbor
dagger call full-pipeline \
  --registry-url=harbor.example.com \
  --registry-username=env:REGISTRY_USER \
  --registry-password=env:REGISTRY_PASSWORD \
  --image-ref=harbor.example.com/myproject/search-api \
  --tag=v1.0.0

# Example 2: GitHub Container Registry (GHCR)
dagger call full-pipeline \
  --registry-url=ghcr.io \
  --registry-username=env:GITHUB_USER \
  --registry-password=env:GITHUB_TOKEN \
  --image-ref=ghcr.io/myorg/search-api \
  --tag=v1.0.0

# Example 3: Docker Hub
dagger call full-pipeline \
  --registry-url=docker.io \
  --registry-username=env:DOCKER_USER \
  --registry-password=env:DOCKER_PASSWORD \
  --image-ref=myusername/search-api \
  --tag=v1.0.0

# Example 4: GitLab Container Registry
dagger call full-pipeline \
  --registry-url=registry.gitlab.com \
  --registry-username=env:GITLAB_USER \
  --registry-password=env:GITLAB_TOKEN \
  --image-ref=registry.gitlab.com/mygroup/myproject/search-api \
  --tag=v1.0.0
```

### Individual Pipeline Steps

```bash
# Security Gates (no --source needed, defaults to current directory)
dagger call secret-scan              # Scan for hardcoded secrets (TruffleHog)
dagger call sast-scan                # Static application security testing (Semgrep)
dagger call dependency-scan          # Dependency vulnerability scan (Trivy)
dagger call license-scan             # License compliance scan (Trivy)
dagger call iac-scan                 # Infrastructure as Code scan (Checkov)
dagger call policy-check             # Policy as Code validation (OPA/Conftest)

# Build and Test
dagger call build                    # Build and run unit tests
dagger call static-analysis          # Code quality checks

# C# Specific Security & Quality
dagger call c-sharp-security-analysis  # .NET analyzers with security rules (enforced)
dagger call code-coverage              # Code coverage with minimum threshold (default 80%)
dagger call code-coverage --minimum-coverage=90  # Custom coverage threshold

# Quality Testing
dagger call mutation-test            # Mutation testing with Stryker.NET (default 80% threshold)
dagger call mutation-test --minimum-score=90  # Custom mutation score threshold

# SBOM and Container
dagger call generate-sbom            # Generate software bill of materials
dagger call build-container          # Build container image
dagger call scan-container \         # Scan container for vulnerabilities
  --container=$(dagger call build-container)

# Supply Chain Security
dagger call sign-image \             # Sign container image with Cosign
  --container=$(dagger call build-container) \
  --private-key=env:COSIGN_PRIVATE_KEY \
  --password=env:COSIGN_PASSWORD \
  --image-ref=harbor.example.com/myproject/search-api:v1.0.0

dagger call attest-sbom \            # Attach signed SBOM attestation
  --sbom="$(dagger call generate-sbom)" \
  --private-key=env:COSIGN_PRIVATE_KEY \
  --password=env:COSIGN_PASSWORD \
  --image-ref=harbor.example.com/myproject/search-api:v1.0.0

dagger call cis-benchmark \          # CIS Docker Benchmark compliance
  --container=$(dagger call build-container)

# Setup K3s cluster for testing
dagger call setup-k3s

# Run integration tests
dagger call run-integration-tests \
  --cluster=$(dagger call setup-k3s)

# Runtime Security & Performance Testing
dagger call dast-scan \              # OWASP ZAP dynamic security testing
  --cluster=$(dagger call setup-k3s)

dagger call api-security-test \      # Nuclei API security testing (OWASP API Top 10)
  --cluster=$(dagger call setup-k3s)

dagger call performance-test \       # k6 load testing
  --cluster=$(dagger call setup-k3s) \
  --virtual-users=50 \
  --duration=2m
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

### Comprehensive Shift-Left Security

This pipeline implements **defense-in-depth** with multiple security layers:

**1. Secret Detection** 🔐
- Tool: TruffleHog
- Scans for hardcoded credentials, API keys, tokens
- Features: Secret verification, 800+ credential detectors
- Detects: AWS keys, GitHub tokens, Slack tokens, database credentials, etc.
- Enforcement: **BLOCKS** pipeline on detection

**2. Static Application Security Testing (SAST)** 🛡️
- Tool: Semgrep
- Detects: SQL injection, XSS, insecure deserialization, crypto issues
- Rulesets: C# security, security-audit, OWASP Top 10
- Enforcement: **BLOCKS** on ERROR/WARNING severity

**3. Dynamic Application Security Testing (DAST)** 🎯
- Tool: OWASP ZAP (Zed Attack Proxy)
- Tests: Running application for vulnerabilities
- Detects: XSS, SQL injection, authentication bypasses, OWASP Top 10
- Method: Baseline scan with spidering and active scanning
- Enforcement: **BLOCKS** on any vulnerability detected

**4. Dependency Vulnerability Scanning** 🔒
- Tool: Trivy (filesystem mode)
- Scans: NuGet packages and transitive dependencies
- Enforcement: **BLOCKS** on HIGH/CRITICAL vulnerabilities

**5. Infrastructure as Code (IaC) Security** ☸️
- Tool: Checkov
- Validates: Kubernetes manifests for misconfigurations
- Checks: Privileged containers, resource limits, RBAC, network policies

**6. Container Security** 🐳
- Tool: Trivy (image mode)
- Scans: OS packages, application dependencies, layers
- Enforcement: **BLOCKS** on HIGH/CRITICAL vulnerabilities

**7. Software Bill of Materials (SBOM)** 📋
- Tool: Syft
- Format: SPDX JSON
- Tracks: All dependencies for supply chain transparency

**8. Runtime Security Hardening** 🔧
- Non-root user execution (searchapi:searchapi)
- Multi-stage builds (minimize attack surface)
- Resource limits (CPU, memory)
- Security contexts in Kubernetes
- Official Microsoft base images only

### C# / .NET Specific Security

**C# Security Analyzers** 🔍
- Tool: Built-in .NET Analyzers
- Analysis Mode: AllEnabledByDefault
- Analysis Level: Latest
- Enforcement: TreatWarningsAsErrors=true
- Detects:
  * Insecure cryptography usage
  * SQL injection vulnerabilities
  * XSS vulnerabilities
  * Insecure deserialization
  * Information disclosure
  * Authentication/authorization issues
  * Regex DoS (ReDoS)
  * And 400+ other .NET specific issues

**Code Coverage** 📊
- Tool: XPlat Code Coverage (built-in)
- Default threshold: 80%
- Format: Cobertura XML
- Enforcement: Configurable minimum coverage
- Tracks: Line, branch, and method coverage

**Benefits of .NET Analyzers:**
- Language-aware analysis (understands C# semantics)
- Catches .NET framework-specific issues
- Enforced at build time
- No external dependencies needed
- Continuously updated by Microsoft

### Advanced Security Features

**License Compliance Scanning** 📜
- Tool: Trivy (license scanner)
- Purpose: Detect problematic licenses in dependencies
- Detects:
  * GPL/AGPL in commercial code
  * License incompatibilities
  * Missing license information
  * Restrictive licenses
- Enforcement: BLOCKS on HIGH/CRITICAL license issues
- Output: JSON report with full license details

**API Security Testing** 🔓
- Tool: Nuclei
- Purpose: Test for OWASP API Security Top 10
- Tests:
  * Broken Object Level Authorization (BOLA)
  * Broken Authentication
  * Excessive Data Exposure
  * Lack of Resources & Rate Limiting
  * Broken Function Level Authorization
  * Mass Assignment
  * Security Misconfiguration
  * Injection
  * Improper Assets Management
  * Insufficient Logging & Monitoring
- Method: Template-based vulnerability scanning
- Enforcement: BLOCKS on HIGH/CRITICAL API vulnerabilities

**Performance Testing** 🚀
- Tool: k6 (Grafana Labs)
- Purpose: Validate API performance under load
- Metrics:
  * Response time (p95 < 500ms)
  * Error rate (< 5%)
  * Throughput
  * Concurrent users
- Configurable:
  * Virtual users (default: 10)
  * Duration (default: 30s)
  * Custom thresholds
- Enforcement: Optional (warns on threshold violations)

**Mutation Testing** 🧬
- Tool: Stryker.NET
- Purpose: Test the quality of your tests
- Method:
  * Mutates source code (changes operators, values, logic)
  * Runs tests against mutated code
  * Verifies tests catch the mutations
- Metrics:
  * Mutation score (% of mutations caught)
  * Survived mutations (tests didn't catch)
  * Killed mutations (tests caught)
- Enforcement: Optional, default 80% threshold (can be slow)

**Image Signing** ✍️
- Tool: Cosign (Sigstore)
- Purpose: Supply chain security and image integrity
- Features:
  * Cryptographic signing of container images
  * Verification of image authenticity
  * Integration with Sigstore transparency log
  * Support for airgapped environments
- Usage: Optional, requires private key and password
- Benefits:
  * Ensures image hasn't been tampered with
  * Proves provenance of the image
  * Meets compliance requirements (e.g., SLSA)

**SBOM Attestation** 📋✍️
- Tool: Cosign (Sigstore)
- Purpose: Cryptographically signed software bill of materials
- Features:
  * Attaches SBOM as in-toto attestation to container image
  * SPDX JSON format predicate
  * Verifiable with cosign verify-attestation
  * Stored in OCI registry alongside image
- Benefits:
  * Immutable dependency tracking
  * Tamper-proof supply chain transparency
  * Compliance with SLSA Level 3
  * Enables automated vulnerability tracking
- Usage: Optional, requires private key and password

**Policy as Code** 📐
- Tool: OPA/Conftest
- Purpose: Validate configurations against custom policies
- Validates:
  * Kubernetes manifests for security requirements
  * Non-root execution enforcement
  * Resource limits (CPU, memory)
  * No privileged containers
  * Custom organizational policies
- Features:
  * Rego policy language (Open Policy Agent)
  * JSON output for CI integration
  * Extensible with custom rules
  * Shift-left policy enforcement
- Enforcement: BLOCKS on policy violations
- Benefits:
  * Consistent security policies
  * Prevent misconfigurations before deployment
  * Self-documenting security requirements

**CIS Benchmark Compliance** 📊
- Tool: Trivy (compliance mode)
- Purpose: Validate Docker containers against CIS Docker Benchmark
- Validates:
  * Image and container configuration
  * Docker security best practices
  * CIS Docker Benchmark v1.6.0
  * Industry-standard security controls
- Checks:
  * User namespaces and privileges
  * Capability restrictions
  * Content trust and verification
  * Network security
  * Logging and auditing
- Output: JSON compliance report with pass/fail status
- Benefits:
  * Industry-recognized security standard
  * Compliance documentation
  * Baseline security validation

### Security Tools Integration

| Category | Tool | Purpose | Enforcement |
|----------|------|---------|-------------|
| Secrets | TruffleHog | Find & verify leaked credentials (800+ detectors) | ✅ Enforced |
| SAST (Generic) | Semgrep | Code vulnerability analysis (OWASP Top 10) | ✅ Enforced |
| SAST (C# Specific) | .NET Analyzers | C#/.NET security issues (400+ rules) | ✅ Enforced |
| Code Coverage | XPlat Coverage | Test coverage enforcement | ✅ Enforced (80%) |
| Mutation Testing | Stryker.NET | Test quality verification | ⚠️ Optional (80%) |
| Dependencies | Trivy | Package vulnerabilities | ✅ Enforced |
| License Compliance | Trivy | License scanning (GPL, AGPL detection) | ✅ Enforced |
| Policy as Code | OPA/Conftest | Custom policy validation (Rego) | ✅ Enforced |
| DAST | OWASP ZAP | Runtime vulnerability testing | ✅ Enforced |
| API Security | Nuclei | OWASP API Security Top 10 | ✅ Enforced |
| Performance | k6 | Load testing & SLA validation | ⚠️ Optional |
| IaC | Checkov | K8s configuration security | ℹ️ Report |
| Container | Trivy | Image vulnerabilities | ✅ Enforced |
| CIS Benchmark | Trivy | Docker CIS compliance (v1.6.0) | ℹ️ Report |
| SBOM | Syft | Dependency tracking (SPDX format) | ℹ️ Generated |
| SBOM Attestation | Cosign | Signed SBOM (in-toto attestation) | ⚠️ Optional |
| Image Signing | Cosign | Supply chain integrity (Sigstore) | ⚠️ Optional |

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
        - --registry-url={{workflow.parameters.registry-url}}
        - --registry-username={{workflow.parameters.registry-username}}
        - --registry-password={{workflow.parameters.registry-password}}
        - --image-ref={{workflow.parameters.image-ref}}
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

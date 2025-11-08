# Dagger CI/CD Pipeline Showcase

## Overview

This repository is a **proof-of-concept demonstration** of building modern CI/CD pipelines with [Dagger](https://dagger.io). The primary goal is to showcase Dagger's capabilities, not to build a production search system.

## What is Dagger?

Dagger is a programmable CI/CD engine that runs your pipelines in containers. Instead of YAML configuration, you write your pipeline as code (in Go, Python, or TypeScript), making it:

- **Portable** - Run locally, in CI, or in Codespaces
- **Reproducible** - Containerized execution ensures consistency
- **Debuggable** - Interactive testing of pipeline steps
- **Fast** - Intelligent caching and parallelization

## This Demonstration

### Example Application: Search API

We chose a C# Search API for Riksarkivet metadata as the example because it demonstrates:

1. **Multi-language support** (.NET/C# in Dagger written in Go)
2. **Complex dependencies** (Solr database, external APIs)
3. **Real-world requirements** (security scanning, testing, deployment)
4. **Infrastructure as Code** (Kubernetes manifests, Helm)

But the **focus is on the pipeline**, not the search functionality.

### The 12-Step Pipeline

Our Dagger pipeline demonstrates:

```go
// .dagger/main.go - Complete CI/CD in Go code

// Step 1-2: Build and Test
dagger call build --source=.
dagger call static-analysis --source=.

// Step 3-4: Security
dagger call security-scan --source=.
dagger call generate-sbom --source=.

// Step 5-6: Container
dagger call build-container --source=.
dagger call scan-container --container=$(...)

// Step 7-11: Deploy & Test
dagger call setup-k3s
dagger call deploy-solr ...
dagger call deploy-api ...
dagger call run-integration-tests ...

// Step 12: Publish
dagger call push-to-harbor ...
```

### Key Patterns Demonstrated

#### 1. **Ephemeral Infrastructure**
```go
// Spin up K3s cluster, deploy services, test, tear down - all in Dagger
func (m *SearchApi) SetupK3s(ctx context.Context) (*K3sCluster, error) {
    k3sContainer := dag.Container().
        From("rancher/k3s:v1.28.5-k3s1").
        WithServiceBinding("registry", registry).
        // ... configure K3s
    return &K3sCluster{Service: k3sService, Kubeconfig: kubeconfig}, nil
}
```

#### 2. **Service Composition**
```go
// Bind services together (registry + K3s + Solr + API)
container.
    WithServiceBinding("k3s", cluster.Service).
    WithServiceBinding("solr", solrService).
    WithServiceBinding("registry", registry)
```

#### 3. **Multi-Registry Publishing**
```go
// Push to different registries with same code
func (m *SearchApi) PushToHarbor(...) (string, error)
func (m *SearchApi) PushToLocalRegistry(...) (string, error)
```

#### 4. **Security Integration**
```go
// SBOM, vulnerability scanning, static analysis - all built in
dagger call generate-sbom --source=.
dagger call scan-container --container=$(...)
```

#### 5. **Integration Testing**
```go
// Run real tests against deployed services in K3s
func (m *SearchApi) RunIntegrationTests(ctx, source, cluster) (string, error)
```

## Running the Demo

### Prerequisites
- Docker or Podman
- Dagger CLI: `curl -L https://dl.dagger.io/dagger/install.sh | sh`

### Quick Demo

```bash
# Clone the repository
git clone <repo-url>
cd search-api

# Run individual steps
dagger call build --source=.
dagger call security-scan --source=.
dagger call build-container --source=.

# Run the complete pipeline
dagger call full-pipeline --source=. --tag=demo

# With Harbor registry push
dagger call full-pipeline \
  --source=. \
  --harbor-url=harbor.example.com \
  --harbor-username=env:HARBOR_USERNAME \
  --harbor-password=env:HARBOR_PASSWORD \
  --harbor-project=search-api \
  --tag=v1.0.0
```

### GitHub Codespaces

The repository includes `.devcontainer` configuration:

1. Click **Code** → **Codespaces** → **Create codespace**
2. Wait for environment to set up (Dagger auto-installs)
3. Run `dagger call full-pipeline --source=.`

Everything works in the browser!

### Local Development

```bash
# Start services
docker-compose up -d

# Run the API
cd SearchApi && dotnet run

# Access Swagger UI
open http://localhost:8080/swagger
```

## Integration Examples

### Argo Workflows

```yaml
# .argo/workflow-template.yaml
- name: run-dagger-pipeline
  container:
    image: dagger:latest
    command: [dagger, call, full-pipeline, --source=/workspace]
```

### GitHub Actions

```yaml
name: CI
on: [push]
jobs:
  dagger:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: dagger/dagger-for-github@v5
      - run: dagger call full-pipeline --source=.
```

### GitLab CI

```yaml
dagger-pipeline:
  image: dagger:latest
  script:
    - dagger call full-pipeline --source=.
```

**Same Dagger code, different CI platforms!**

## Key Takeaways

### What This Demo Shows

✅ **Portable Pipelines** - Run anywhere Docker runs
✅ **Developer-Friendly** - Debug locally, same as CI
✅ **Type-Safe** - Go functions, not YAML strings
✅ **Fast** - Intelligent caching, parallel execution
✅ **Secure** - Built-in security scanning, SBOM generation
✅ **Testable** - Integration tests in ephemeral K3s clusters
✅ **Flexible** - Easy to extend and modify

### What This Is NOT

❌ A production search API implementation
❌ A comprehensive Solr guide
❌ A Riksarkivet integration tutorial
❌ A Kubernetes best practices example

Those are all secondary. **This is about Dagger.**

## Resources

- **Dagger Documentation**: https://docs.dagger.io
- **Dagger Examples**: https://github.com/dagger/dagger/tree/main/examples
- **This Pipeline**: `.dagger/main.go`
- **Inspired By**: [AI-Riksarkivet/coder-templates](https://github.com/AI-Riksarkivet/coder-templates)

## Questions?

- **"Why Go for the pipeline?"** - Dagger supports Go, Python, TypeScript. Go was chosen for type safety and performance.
- **"Can I use this for my project?"** - Yes! The `.dagger/main.go` is a template you can adapt.
- **"Does it work on my CI platform?"** - Yes! Dagger works on GitHub Actions, GitLab, Jenkins, Argo, CircleCI, etc.
- **"How fast is it?"** - With caching, rebuilds are typically < 1 minute for code-only changes.

## Next Steps

Want to adapt this for your project?

1. **Copy `.dagger/main.go`** as a starting point
2. **Replace the C# build** with your language (Go, Python, Node.js, etc.)
3. **Customize the pipeline steps** for your needs
4. **Keep the patterns**: security scanning, SBOM, testing, deployment
5. **Run it**: `dagger call full-pipeline --source=.`

The whole point of Dagger is **you write the pipeline once, run it everywhere**.

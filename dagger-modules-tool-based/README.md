# Tool-Based Dagger Modules

**Language-agnostic, reusable security and CI/CD modules for Daggerverse**

## 🎯 Philosophy

Each module wraps **one CLI tool**, making them:
- ✅ **Language Agnostic** - Works with Python, Go, Java, .NET, JavaScript, etc.
- ✅ **Simple** - One tool = one module
- ✅ **Composable** - Mix and match for your pipeline
- ✅ **Maintainable** - Tool updates only affect one module
- ✅ **Community Friendly** - Anyone can use and contribute

## 📦 Available Modules

### Security Scanning
| Module | Tool | Purpose | Languages |
|--------|------|---------|-----------|
| [trufflehog](./trufflehog/) | TruffleHog | Secret scanning | All |
| [semgrep](./semgrep/) | Semgrep | SAST (static analysis) | 30+ languages |
| [trivy](./trivy/) | Trivy | Vulnerabilities, licenses, secrets, misconfigs | All |

### Dynamic Testing
| Module | Tool | Purpose | Type |
|--------|------|---------|------|
| [zap](./zap/) | OWASP ZAP | DAST (dynamic app scanning) | Web apps |
| [nuclei](./nuclei/) | Nuclei | Template-based security testing | APIs, web apps |
| [k6](./k6/) | k6 | Load & performance testing | APIs, web apps |

### Supply Chain Security
| Module | Tool | Purpose | Works With |
|--------|------|---------|------------|
| [syft](./syft/) | Syft | SBOM generation | Source code, containers, images |
| [cosign](./cosign/) | Cosign | Image signing & verification | Container images |

### Infrastructure as Code
| Module | Tool | Purpose | Supports |
|--------|------|---------|----------|
| [checkov](./checkov/) | Checkov | IaC security scanning | K8s, Terraform, CloudFormation, Dockerfile |
| [conftest](./conftest/) | Conftest | OPA policy validation | K8s, Terraform, any config files |

## 🚀 Quick Start

### Using a Single Module

```bash
# Scan for secrets (works with ANY language)
dagger call -m ./dagger-modules-tool-based/trufflehog scan --source=.

# Run SAST on Python code
dagger call -m ./dagger-modules-tool-based/semgrep scan-language \
  --source=. \
  --language=python

# Scan container for vulnerabilities
dagger call -m ./dagger-modules-tool-based/trivy scan-container \
  --container=<container>
```

### Composing Multiple Modules

Create a pipeline that uses multiple tools:

```bash
# 1. Secrets
dagger call -m ./dagger-modules-tool-based/trufflehog scan --source=.

# 2. SAST
dagger call -m ./dagger-modules-tool-based/semgrep scan-language \
  --source=. \
  --language=go

# 3. Dependencies
dagger call -m ./dagger-modules-tool-based/trivy scan-vulnerabilities \
  --source=.

# 4. Licenses
dagger call -m ./dagger-modules-tool-based/trivy scan-licenses \
  --source=.

# 5. IaC
dagger call -m ./dagger-modules-tool-based/checkov scan-kubernetes \
  --source=. \
  --k8s-dir=k8s
```

## 💡 Why Tool-Based?

### ❌ Old Approach (Language-Specific)
```
dotnet-security/    # Only works with .NET
python-security/    # Only works with Python
go-security/        # Only works with Go
java-security/      # Only works with Java
```

**Problems:**
- Duplicate effort for each language
- Not reusable across tech stacks
- Harder to maintain

### ✅ New Approach (Tool-Based)
```
trufflehog/         # Works with ALL languages
semgrep/            # Works with 30+ languages
trivy/              # Works with ALL languages
```

**Benefits:**
- Write once, use everywhere
- Community can use for any project
- Each module is simple and focused
- Easy to maintain and upgrade

## 📚 Module Details

### 1. trufflehog - Secret Scanning
Scans for hardcoded secrets (API keys, passwords, tokens).

**Works with:** Any programming language, Git repos, Docker images

**Example:**
```bash
# Scan filesystem
dagger call -m ./dagger-modules-tool-based/trufflehog scan --source=.

# Scan Git history
dagger call -m ./dagger-modules-tool-based/trufflehog scan-git \
  --repo-url=https://github.com/user/repo

# Scan Docker image
dagger call -m ./dagger-modules-tool-based/trufflehog scan-docker \
  --container=<container>
```

---

### 2. semgrep - SAST Scanner
Static Application Security Testing for 30+ languages.

**Supports:** Python, JavaScript, Go, Java, C#, Ruby, PHP, Rust, and more

**Example:**
```bash
# Auto-detect language
dagger call -m ./dagger-modules-tool-based/semgrep scan --source=.

# Specific language with OWASP rules
dagger call -m ./dagger-modules-tool-based/semgrep scan-language \
  --source=. \
  --language=python \
  --security-audit=true \
  --owasp-top-ten=true

# Scan for XSS vulnerabilities
dagger call -m ./dagger-modules-tool-based/semgrep scan-xss --source=.
```

---

### 3. trivy - Comprehensive Scanner
Scans for vulnerabilities, licenses, secrets, and misconfigurations.

**Scans:** Dependencies, containers, IaC, licenses, secrets

**Example:**
```bash
# Scan for dependency vulnerabilities
dagger call -m ./dagger-modules-tool-based/trivy scan-vulnerabilities \
  --source=. \
  --severity=HIGH,CRITICAL

# Scan for problematic licenses
dagger call -m ./dagger-modules-tool-based/trivy scan-licenses \
  --source=.

# Scan container image
dagger call -m ./dagger-modules-tool-based/trivy scan-container \
  --container=<container>

# Scan Kubernetes manifests
dagger call -m ./dagger-modules-tool-based/trivy scan-kubernetes \
  --source=.
```

---

### 4. syft - SBOM Generator
Generates Software Bill of Materials (SBOM).

**Formats:** SPDX-JSON, CycloneDX-JSON, Syft-JSON

**Example:**
```bash
# Generate SBOM from source
dagger call -m ./dagger-modules-tool-based/syft scan \
  --source=. \
  --format=spdx-json

# Generate SBOM from container
dagger call -m ./dagger-modules-tool-based/syft scan-container \
  --container=<container>
```

---

### 5. cosign - Image Signing
Signs and verifies container images with Sigstore Cosign.

**Example:**
```bash
# Sign an image
dagger call -m ./dagger-modules-tool-based/cosign sign \
  --container=<container> \
  --private-key=env:COSIGN_KEY \
  --password=env:COSIGN_PASSWORD \
  --image-ref="myregistry.com/app:v1.0"

# Verify an image
dagger call -m ./dagger-modules-tool-based/cosign verify \
  --image-ref="myregistry.com/app:v1.0" \
  --public-key=env:COSIGN_PUBLIC_KEY
```

---

### 6. zap - DAST Scanner
OWASP ZAP for dynamic application security testing.

**Example:**
```bash
# Baseline scan (quick, passive)
dagger call -m ./dagger-modules-tool-based/zap baseline-scan \
  --api-service=<service> \
  --target-url="http://api:8080"

# Full scan (comprehensive, active)
dagger call -m ./dagger-modules-tool-based/zap full-scan \
  --api-service=<service> \
  --target-url="http://api:8080"
```

---

### 7. nuclei - Security Testing
Template-based vulnerability scanner.

**Example:**
```bash
# API security scan
dagger call -m ./dagger-modules-tool-based/nuclei scan-api \
  --api-service=<service> \
  --target-url="http://api:8080"

# CVE scan
dagger call -m ./dagger-modules-tool-based/nuclei scan-cve \
  --api-service=<service> \
  --target-url="http://api:8080"
```

---

### 8. k6 - Load Testing
Performance and load testing.

**Example:**
```bash
# Simple load test
dagger call -m ./dagger-modules-tool-based/k6 load-test \
  --api-service=<service> \
  --target-url="http://api:8080" \
  --endpoint="/api/search" \
  --vus=50 \
  --duration="2m"

# Stress test with ramping
dagger call -m ./dagger-modules-tool-based/k6 stress-test \
  --api-service=<service> \
  --max-vus=100
```

---

### 9. checkov - IaC Scanner
Scans Infrastructure as Code for security issues.

**Supports:** Kubernetes, Terraform, CloudFormation, ARM templates, Dockerfiles, Helm

**Example:**
```bash
# Scan Kubernetes manifests
dagger call -m ./dagger-modules-tool-based/checkov scan-kubernetes \
  --source=. \
  --k8s-dir=k8s

# Scan Terraform
dagger call -m ./dagger-modules-tool-based/checkov scan-terraform \
  --source=. \
  --terraform-dir=terraform
```

---

### 10. conftest - Policy Validation
OPA (Open Policy Agent) policy testing.

**Example:**
```bash
# Test Kubernetes manifests against policies
dagger call -m ./dagger-modules-tool-based/conftest test-kubernetes \
  --source=. \
  --k8s-dir=k8s

# Test with custom policies
dagger call -m ./dagger-modules-tool-based/conftest test \
  --source=. \
  --input=k8s \
  --policy-dir=./policies
```

---

## 🔄 Complete Security Pipeline Example

Here's how to compose all modules into a complete pipeline:

```bash
#!/bin/bash
set -e

echo "🔐 Phase 1: Static Security Analysis"
dagger call -m ./dagger-modules-tool-based/trufflehog scan --source=.
dagger call -m ./dagger-modules-tool-based/semgrep scan-language --source=. --language=python
dagger call -m ./dagger-modules-tool-based/trivy scan-vulnerabilities --source=.
dagger call -m ./dagger-modules-tool-based/trivy scan-licenses --source=.

echo "☸️ Phase 2: IaC Security"
dagger call -m ./dagger-modules-tool-based/checkov scan-kubernetes --source=.
dagger call -m ./dagger-modules-tool-based/conftest test-kubernetes --source=.

echo "📦 Phase 3: Build Container"
# (your container build here)

echo "🔎 Phase 4: Container Security"
dagger call -m ./dagger-modules-tool-based/trivy scan-container --container=$CONTAINER

echo "📋 Phase 5: Supply Chain"
dagger call -m ./dagger-modules-tool-based/syft scan-container --container=$CONTAINER
dagger call -m ./dagger-modules-tool-based/cosign sign --container=$CONTAINER

echo "🎯 Phase 6: Dynamic Testing"
# (start your service)
dagger call -m ./dagger-modules-tool-based/zap baseline-scan --api-service=$SERVICE
dagger call -m ./dagger-modules-tool-based/nuclei scan-api --api-service=$SERVICE
dagger call -m ./dagger-modules-tool-based/k6 load-test --api-service=$SERVICE

echo "✅ All security checks passed!"
```

## 📖 Publishing to Daggerverse

Each module can be published independently:

```bash
cd trufflehog
git init
git add .
git commit -m "Initial commit"
gh repo create yourorg/dagger-trufflehog --public --source=. --remote=origin
git push -u origin main
dagger publish
```

Once published, anyone can use:

```bash
dagger install github.com/yourorg/dagger-trufflehog
dagger call -m github.com/yourorg/dagger-trufflehog scan --source=.
```

## 🤝 Contributing

These modules are designed to be community-driven. Contributions welcome!

See [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

## 📝 License

Apache-2.0 (see LICENSE in each module directory)

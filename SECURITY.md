# Security Policy

## 🔒 Security First Approach

This project implements a **shift-left security** strategy, integrating security at every stage of the development lifecycle.

## Security Features

### 1. Build-Time Security

- **Dependency Scanning**: Automated scanning of NuGet packages for known vulnerabilities
- **Static Analysis**: Code quality and security pattern detection
- **SBOM Generation**: Complete software bill of materials for transparency
- **Container Scanning**: Multi-layer vulnerability analysis with Trivy

### 2. Runtime Security

- **Non-Root Execution**: All containers run as unprivileged users
- **Resource Limits**: CPU and memory constraints prevent resource exhaustion
- **Health Checks**: Automated health monitoring and restart policies
- **Least Privilege**: Minimal permissions and capabilities

### 3. Infrastructure Security

- **Network Policies**: Kubernetes network segmentation
- **Secret Management**: Integration with Infisical for secure secret storage
- **Registry Security**: Private Harbor registry with vulnerability scanning
- **TLS/HTTPS**: Encrypted communication channels

## CI/CD Security Pipeline

Our Dagger pipeline includes the following security steps:

1. **Dependency Vulnerability Scan**
   ```bash
   dagger call security-scan --source=.
   ```

2. **Static Code Analysis**
   ```bash
   dagger call static-analysis --source=.
   ```

3. **SBOM Generation**
   ```bash
   dagger call generate-sbom --source=.
   ```

4. **Container Security Scan**
   ```bash
   dagger call scan-container --container=$(dagger call build-container --source=.)
   ```

## Reporting a Vulnerability

If you discover a security vulnerability, please:

1. **Do NOT** open a public issue
2. Email security@your-domain.com with:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

We will respond within 48 hours and work with you to address the issue.

## Security Best Practices for Contributors

1. **Never commit secrets** - Use environment variables or secret management
2. **Keep dependencies updated** - Regularly update NuGet packages
3. **Follow secure coding practices** - Input validation, output encoding
4. **Run security scans** - Before submitting PRs
5. **Review security policies** - Stay informed about security requirements

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Security Tools Used

- **Trivy**: Container vulnerability scanning
- **Syft**: SBOM generation
- **dotnet list package --vulnerable**: Dependency scanning
- **Infisical**: Secret management
- **Harbor**: Secure container registry with scanning

## Compliance

This project follows:

- OWASP Top 10 security practices
- CIS Docker Benchmarks
- Kubernetes Pod Security Standards

## Security Updates

Security patches are released as soon as possible after discovery. Subscribe to releases to stay informed.

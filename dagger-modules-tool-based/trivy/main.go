// Dagger module for Trivy - comprehensive vulnerability scanner
// Scans: vulnerabilities, misconfigurations, secrets, licenses, SBOMs
package main

import (
	"context"
	"dagger/trivy/internal/dagger"
)

type Trivy struct{}

// ScanFilesystem scans source code for vulnerabilities, secrets, misconfigs, licenses
func (m *Trivy) ScanFilesystem(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Scanners to use: vuln, misconfig, secret, license
	// +default=["vuln"]
	scanners []string,
	// Severity levels: LOW, MEDIUM, HIGH, CRITICAL
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Output format: json, table, sarif, cyclonedx, spdx, github
	// +default="json"
	format string,
	// Exit code when vulnerabilities are found (0 = no fail, 1 = fail)
	// +default=0
	exitCode int,
) (string, error) {
	scannersStr := ""
	for i, s := range scanners {
		if i > 0 {
			scannersStr += ","
		}
		scannersStr += s
	}

	severityStr := ""
	for i, s := range severity {
		if i > 0 {
			severityStr += ","
		}
		severityStr += s
	}

	args := []string{
		"trivy", "fs",
		"--scanners", scannersStr,
		"--severity", severityStr,
		"--format", format,
	}

	if exitCode > 0 {
		args = append(args, "--exit-code", "1")
	}

	args = append(args, ".")

	return dag.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

// ScanContainer scans a container image for vulnerabilities
func (m *Trivy) ScanContainer(
	ctx context.Context,
	// Container to scan
	container *dagger.Container,
	// Scanners to use
	// +default=["vuln"]
	scanners []string,
	// Severity levels
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Output format
	// +default="json"
	format string,
	// Exit code on findings
	// +default=0
	exitCode int,
) (string, error) {
	tarball := container.AsTarball()

	scannersStr := ""
	for i, s := range scanners {
		if i > 0 {
			scannersStr += ","
		}
		scannersStr += s
	}

	severityStr := ""
	for i, s := range severity {
		if i > 0 {
			severityStr += ","
		}
		severityStr += s
	}

	args := []string{
		"trivy", "image",
		"--input", "/image.tar",
		"--scanners", scannersStr,
		"--severity", severityStr,
		"--format", format,
	}

	if exitCode > 0 {
		args = append(args, "--exit-code", "1")
	}

	return dag.Container().
		From("aquasec/trivy:latest").
		WithMountedFile("/image.tar", tarball).
		WithExec(args).
		Stdout(ctx)
}

// ScanVulnerabilities scans for package vulnerabilities (dependencies)
func (m *Trivy) ScanVulnerabilities(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Severity levels
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Fail build on findings
	// +default=true
	failOnFindings bool,
) (string, error) {
	exitCode := 0
	if failOnFindings {
		exitCode = 1
	}

	return m.ScanFilesystem(ctx, source, []string{"vuln"}, severity, "json", exitCode)
}

// ScanLicenses scans for license compliance issues
func (m *Trivy) ScanLicenses(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Severity levels for problematic licenses
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Fail build on problematic licenses
	// +default=true
	failOnFindings bool,
) (string, error) {
	exitCode := 0
	if failOnFindings {
		exitCode = 1
	}

	return m.ScanFilesystem(ctx, source, []string{"license"}, severity, "json", exitCode)
}

// ScanSecrets scans for hardcoded secrets in source code
func (m *Trivy) ScanSecrets(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Fail build on secrets found
	// +default=true
	failOnFindings bool,
) (string, error) {
	exitCode := 0
	if failOnFindings {
		exitCode = 1
	}

	return m.ScanFilesystem(ctx, source, []string{"secret"}, []string{"HIGH", "CRITICAL"}, "json", exitCode)
}

// ScanMisconfigs scans for IaC misconfigurations (Kubernetes, Terraform, Docker, etc.)
func (m *Trivy) ScanMisconfigs(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Severity levels
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Fail build on misconfigurations
	// +default=false
	failOnFindings bool,
) (string, error) {
	exitCode := 0
	if failOnFindings {
		exitCode = 1
	}

	return m.ScanFilesystem(ctx, source, []string{"misconfig"}, severity, "json", exitCode)
}

// ScanAll runs all Trivy scanners (vulnerabilities, secrets, misconfigs, licenses)
func (m *Trivy) ScanAll(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Severity levels
	// +default=["HIGH", "CRITICAL"]
	severity []string,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	return m.ScanFilesystem(
		ctx,
		source,
		[]string{"vuln", "secret", "misconfig", "license"},
		severity,
		format,
		0, // Don't fail, just report
	)
}

// GenerateSbom generates an SBOM (Software Bill of Materials) using Trivy
func (m *Trivy) GenerateSbom(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// SBOM format: cyclonedx, spdx, spdx-json, github
	// +default="spdx-json"
	format string,
) (string, error) {
	return dag.Container().
		From("aquasec/trivy:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"trivy", "fs",
			"--format", format,
			".",
		}).
		Stdout(ctx)
}

// ScanKubernetes scans Kubernetes manifests for security issues
func (m *Trivy) ScanKubernetes(
	ctx context.Context,
	// Source directory containing K8s manifests
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Severity levels
	// +default=["HIGH", "CRITICAL"]
	severity []string,
) (string, error) {
	return m.ScanMisconfigs(ctx, source, severity, false)
}

// Dagger module for Checkov - Infrastructure as Code security scanner
// Supports: Kubernetes, Terraform, CloudFormation, ARM, Dockerfile, and more
package main

import (
	"context"
	"dagger/checkov/internal/dagger"
)

type Checkov struct{}

// Scan runs Checkov on Infrastructure as Code files
func (m *Checkov) Scan(
	ctx context.Context,
	// Source directory containing IaC files
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Framework to scan: kubernetes, terraform, cloudformation, arm, dockerfile, all
	// +default=["all"]
	framework []string,
	// Directory to scan (relative to source)
	// +default="."
	directory string,
	// Fail on severity: critical, high, medium, low
	// +optional
	failOn string,
	// Skip checks (comma-separated check IDs)
	// +optional
	skipChecks []string,
) (string, error) {
	args := []string{"checkov", "-d", directory}

	// Add frameworks
	for _, fw := range framework {
		args = append(args, "--framework", fw)
	}

	// Add fail-on
	if failOn != "" {
		args = append(args, "--check", failOn)
	}

	// Add skip checks
	for _, skip := range skipChecks {
		args = append(args, "--skip-check", skip)
	}

	args = append(args, "--compact", "--quiet")

	return dag.Container().
		From("bridgecrew/checkov:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

// ScanKubernetes scans Kubernetes manifests
func (m *Checkov) ScanKubernetes(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory containing K8s manifests
	// +default="k8s"
	k8sDir string,
) (string, error) {
	return m.Scan(ctx, source, []string{"kubernetes"}, k8sDir, "", nil)
}

// ScanTerraform scans Terraform configurations
func (m *Checkov) ScanTerraform(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory containing Terraform files
	// +default="terraform"
	terraformDir string,
) (string, error) {
	return m.Scan(ctx, source, []string{"terraform"}, terraformDir, "", nil)
}

// ScanDockerfile scans Dockerfiles for security issues
func (m *Checkov) ScanDockerfile(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
) (string, error) {
	return m.Scan(ctx, source, []string{"dockerfile"}, ".", "", nil)
}

// ScanHelm scans Helm charts
func (m *Checkov) ScanHelm(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory containing Helm charts
	// +default="helm"
	helmDir string,
) (string, error) {
	return m.Scan(ctx, source, []string{"helm"}, helmDir, "", nil)
}

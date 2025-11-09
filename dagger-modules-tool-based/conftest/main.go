// Dagger module for Conftest - policy testing using Open Policy Agent (OPA)
package main

import (
	"context"
	"dagger/conftest/internal/dagger"
)

type Conftest struct{}

// Test runs Conftest policy tests on configuration files
func (m *Conftest) Test(
	ctx context.Context,
	// Source directory containing files to test
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory or file to test
	// +default="."
	input string,
	// Directory containing Rego policy files
	// +optional
	policyDir *dagger.Directory,
	// Output format: json, tap, table, junit
	// +default="json"
	outputFormat string,
	// Namespace to use
	// +default="main"
	namespace string,
) (string, error) {
	container := dag.Container().
		From("openpolicyagent/conftest:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src")

	// Use custom policies or create default
	if policyDir != nil {
		container = container.WithDirectory("/policy", policyDir)
	} else {
		// Create a default policy
		defaultPolicy := `package main

deny contains msg if {
  input.kind == "Deployment"
  not input.spec.template.spec.securityContext.runAsNonRoot
  msg := "Containers must not run as root"
}

deny contains msg if {
  input.kind == "Deployment"
  container := input.spec.template.spec.containers[_]
  not container.resources.limits.memory
  msg := sprintf("Container %s must have memory limits", [container.name])
}

deny contains msg if {
  input.kind == "Deployment"
  container := input.spec.template.spec.containers[_]
  not container.resources.limits.cpu
  msg := sprintf("Container %s must have CPU limits", [container.name])
}

deny contains msg if {
  input.kind == "Deployment"
  container := input.spec.template.spec.containers[_]
  container.securityContext.privileged == true
  msg := sprintf("Container %s must not run in privileged mode", [container.name])
}`
		container = container.
			WithExec([]string{"sh", "-c", "mkdir -p /policy"}).
			WithNewFile("/policy/policy.rego", defaultPolicy)
	}

	args := []string{
		"conftest", "test",
		input,
		"--policy", "/policy",
		"--output", outputFormat,
		"--namespace", namespace,
	}

	return container.WithExec(args).Stdout(ctx)
}

// TestKubernetes tests Kubernetes manifests against policies
func (m *Conftest) TestKubernetes(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory containing K8s manifests
	// +default="k8s"
	k8sDir string,
	// Custom policy directory (optional)
	// +optional
	policyDir *dagger.Directory,
) (string, error) {
	return m.Test(ctx, source, k8sDir, policyDir, "json", "main")
}

// TestDockerfile tests Dockerfiles against policies
func (m *Conftest) TestDockerfile(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Dockerfile path
	// +default="Dockerfile"
	dockerfile string,
	// Custom policy directory (optional)
	// +optional
	policyDir *dagger.Directory,
) (string, error) {
	return m.Test(ctx, source, dockerfile, policyDir, "json", "main")
}

// TestTerraform tests Terraform configurations against policies
func (m *Conftest) TestTerraform(
	ctx context.Context,
	// Source directory
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Terraform directory
	// +default="terraform"
	terraformDir string,
	// Custom policy directory (optional)
	// +optional
	policyDir *dagger.Directory,
) (string, error) {
	return m.Test(ctx, source, terraformDir, policyDir, "json", "main")
}

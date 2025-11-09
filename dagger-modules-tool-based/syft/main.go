// Dagger module for Syft - SBOM (Software Bill of Materials) generator
package main

import (
	"context"
	"dagger/syft/internal/dagger"
)

type Syft struct{}

// Scan generates an SBOM from source code (works with any language)
func (m *Syft) Scan(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Output format: spdx-json, cyclonedx-json, syft-json, table, text
	// +default="spdx-json"
	format string,
) (string, error) {
	return dag.Container().
		From("anchore/syft:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"syft", "scan", ".", "-o", format,
		}).
		Stdout(ctx)
}

// ScanContainer generates an SBOM from a container image
func (m *Syft) ScanContainer(
	ctx context.Context,
	// Container to scan
	container *dagger.Container,
	// Output format
	// +default="spdx-json"
	format string,
) (string, error) {
	tarball := container.AsTarball()

	return dag.Container().
		From("anchore/syft:latest").
		WithMountedFile("/image.tar", tarball).
		WithExec([]string{
			"syft", "scan", "docker-archive:/image.tar", "-o", format,
		}).
		Stdout(ctx)
}

// ScanImage generates an SBOM from a remote container image
func (m *Syft) ScanImage(
	ctx context.Context,
	// Image reference (e.g., "alpine:latest", "myregistry.com/app:v1.0")
	imageRef string,
	// Output format
	// +default="spdx-json"
	format string,
	// Registry credentials (optional)
	// +optional
	registryUsername string,
	// +optional
	registryPassword *dagger.Secret,
) (string, error) {
	container := dag.Container().From("anchore/syft:latest")

	if registryPassword != nil {
		container = container.WithSecretVariable("REGISTRY_PASSWORD", registryPassword)
		container = container.WithEnvVariable("REGISTRY_USERNAME", registryUsername)
	}

	return container.
		WithExec([]string{
			"syft", "scan", imageRef, "-o", format,
		}).
		Stdout(ctx)
}

// ScanGit generates an SBOM from a Git repository
func (m *Syft) ScanGit(
	ctx context.Context,
	// Git repository URL
	repoUrl string,
	// Output format
	// +default="spdx-json"
	format string,
) (string, error) {
	return dag.Container().
		From("anchore/syft:latest").
		WithExec([]string{
			"syft", "scan", repoUrl, "-o", format,
		}).
		Stdout(ctx)
}

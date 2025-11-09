// Dagger module for Dive - container image analysis tool
// Provides layer-by-layer breakdown and efficiency analysis
package main

import (
	"context"
	"dagger/dive/internal/dagger"
)

type Dive struct{}

// Analyze analyzes a container image for size and efficiency
func (m *Dive) Analyze(
	ctx context.Context,
	// Container to analyze
	container *dagger.Container,
	// CI mode for machine-readable output
	// +default=true
	ciMode bool,
	// Source type
	// +default="docker-archive"
	sourceType string,
) (string, error) {
	// Save container as tarball
	tarball := container.AsTarball()

	args := []string{"dive"}

	if sourceType != "" {
		args = append(args, "--source", sourceType)
	}

	if ciMode {
		args = append(args, "--ci")
	}

	args = append(args, "/image.tar")

	return dag.Container().
		From("wagoodman/dive:latest").
		WithMountedFile("/image.tar", tarball).
		WithExec(args).
		Stdout(ctx)
}

// GetSize gets the size of a container image
func (m *Dive) GetSize(
	ctx context.Context,
	// Container to analyze
	container *dagger.Container,
) (string, error) {
	// Save container as tarball
	tarball := container.AsTarball()

	return dag.Container().
		From("alpine:latest").
		WithExec([]string{"apk", "add", "--no-cache", "file"}).
		WithMountedFile("/image.tar", tarball).
		WithExec([]string{"sh", "-c", "ls -lh /image.tar | awk '{print $5}'"}).
		Stdout(ctx)
}

// CompareImages compares two container images for size differences
func (m *Dive) CompareImages(
	ctx context.Context,
	// First container
	container1 *dagger.Container,
	// Second container
	container2 *dagger.Container,
) (string, error) {
	size1, err := m.GetSize(ctx, container1)
	if err != nil {
		return "", err
	}

	size2, err := m.GetSize(ctx, container2)
	if err != nil {
		return "", err
	}

	return "Container 1 size: " + size1 + "\nContainer 2 size: " + size2, nil
}

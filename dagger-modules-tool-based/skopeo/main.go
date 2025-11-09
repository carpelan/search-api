// Dagger module for Skopeo - container image operations
// Provides copy, inspect, delete operations for container images
package main

import (
	"context"
	"dagger/skopeo/internal/dagger"
	"fmt"
)

type Skopeo struct{}

// Copy copies a container image from source to destination
func (m *Skopeo) Copy(
	ctx context.Context,
	// Container to copy
	container *dagger.Container,
	// Destination image reference (e.g., "docker://registry:5000/image:tag")
	destRef string,
	// Service binding for registry (optional)
	// +optional
	registryService *dagger.Service,
	// Disable TLS verification
	// +default=false
	disableTLS bool,
	// Source type
	// +default="docker-archive"
	sourceType string,
) (string, error) {
	// Save container as tarball
	tarball := container.AsTarball()

	args := []string{"skopeo", "copy"}

	if disableTLS {
		args = append(args, "--dest-tls-verify=false")
	}

	args = append(args, fmt.Sprintf("%s:/image.tar", sourceType), destRef)

	c := dag.Container().
		From("quay.io/skopeo/stable:latest").
		WithMountedFile("/image.tar", tarball)

	if registryService != nil {
		c = c.WithServiceBinding("registry", registryService)
	}

	return c.WithExec(args).Stdout(ctx)
}

// Inspect inspects a container image
func (m *Skopeo) Inspect(
	ctx context.Context,
	// Image reference to inspect
	imageRef string,
	// Service binding for registry (optional)
	// +optional
	registryService *dagger.Service,
	// Disable TLS verification
	// +default=false
	disableTLS bool,
) (string, error) {
	args := []string{"skopeo", "inspect"}

	if disableTLS {
		args = append(args, "--tls-verify=false")
	}

	args = append(args, imageRef)

	c := dag.Container().
		From("quay.io/skopeo/stable:latest")

	if registryService != nil {
		c = c.WithServiceBinding("registry", registryService)
	}

	return c.WithExec(args).Stdout(ctx)
}

// Delete deletes an image from a registry
func (m *Skopeo) Delete(
	ctx context.Context,
	// Image reference to delete
	imageRef string,
	// Service binding for registry (optional)
	// +optional
	registryService *dagger.Service,
	// Disable TLS verification
	// +default=false
	disableTLS bool,
) (string, error) {
	args := []string{"skopeo", "delete"}

	if disableTLS {
		args = append(args, "--tls-verify=false")
	}

	args = append(args, imageRef)

	c := dag.Container().
		From("quay.io/skopeo/stable:latest")

	if registryService != nil {
		c = c.WithServiceBinding("registry", registryService)
	}

	return c.WithExec(args).Stdout(ctx)
}

// PushToRegistry pushes a container to a registry
func (m *Skopeo) PushToRegistry(
	ctx context.Context,
	// Container to push
	container *dagger.Container,
	// Registry host (e.g., "registry:5000")
	registryHost string,
	// Image name (e.g., "myapp")
	imageName string,
	// Image tag
	// +default="latest"
	tag string,
	// Registry service binding (optional)
	// +optional
	registryService *dagger.Service,
	// Disable TLS verification
	// +default=false
	disableTLS bool,
) (string, error) {
	destRef := fmt.Sprintf("docker://%s/%s:%s", registryHost, imageName, tag)
	return m.Copy(ctx, container, destRef, registryService, disableTLS, "docker-archive")
}

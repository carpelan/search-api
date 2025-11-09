// Dagger module for Cosign - container image signing and verification
package main

import (
	"context"
	"dagger/cosign/internal/dagger"
)

type Cosign struct{}

// Sign signs a container image with Cosign
func (m *Cosign) Sign(
	ctx context.Context,
	// Container to sign
	container *dagger.Container,
	// Private key for signing
	privateKey *dagger.Secret,
	// Password for the private key
	password *dagger.Secret,
	// Image reference to sign (e.g., "myregistry.com/app:v1.0")
	imageRef string,
	// Upload to transparency log (Rekor)
	// +default=false
	tlogUpload bool,
) (string, error) {
	tarball := container.AsTarball()

	tlogFlag := "--tlog-upload=false"
	if tlogUpload {
		tlogFlag = "--tlog-upload=true"
	}

	return dag.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithMountedFile("/image.tar", tarball).
		WithMountedSecret("/cosign.key", privateKey).
		WithSecretVariable("COSIGN_PASSWORD", password).
		WithExec([]string{
			"cosign", "sign",
			"--key", "/cosign.key",
			tlogFlag,
			imageRef,
		}).
		Stdout(ctx)
}

// Verify verifies a signed container image
func (m *Cosign) Verify(
	ctx context.Context,
	// Image reference to verify
	imageRef string,
	// Public key for verification
	publicKey *dagger.Secret,
) (string, error) {
	return dag.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithMountedSecret("/cosign.pub", publicKey).
		WithExec([]string{
			"cosign", "verify",
			"--key", "/cosign.pub",
			imageRef,
		}).
		Stdout(ctx)
}

// Attest attaches an attestation to a container image
func (m *Cosign) Attest(
	ctx context.Context,
	// Attestation data (e.g., SBOM, provenance)
	attestation string,
	// Private key for signing
	privateKey *dagger.Secret,
	// Password for the private key
	password *dagger.Secret,
	// Image reference to attest
	imageRef string,
	// Predicate type (spdxjson, cyclonedx, slsaprovenance, custom)
	// +default="spdxjson"
	predicateType string,
	// Upload to transparency log
	// +default=false
	tlogUpload bool,
) (string, error) {
	tlogFlag := "--tlog-upload=false"
	if tlogUpload {
		tlogFlag = "--tlog-upload=true"
	}

	return dag.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithNewFile("/attestation.json", attestation).
		WithMountedSecret("/cosign.key", privateKey).
		WithSecretVariable("COSIGN_PASSWORD", password).
		WithExec([]string{
			"cosign", "attest",
			"--key", "/cosign.key",
			"--predicate", "/attestation.json",
			"--type", predicateType,
			tlogFlag,
			imageRef,
		}).
		Stdout(ctx)
}

// GenerateKeyPair generates a Cosign key pair for signing
func (m *Cosign) GenerateKeyPair(
	ctx context.Context,
	// Password for the private key
	password *dagger.Secret,
) (*dagger.Directory, error) {
	return dag.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithSecretVariable("COSIGN_PASSWORD", password).
		WithExec([]string{
			"cosign", "generate-key-pair",
		}).
		Directory("/"), nil
}

// VerifyAttestation verifies an attestation on a container image
func (m *Cosign) VerifyAttestation(
	ctx context.Context,
	// Image reference to verify
	imageRef string,
	// Public key for verification
	publicKey *dagger.Secret,
	// Predicate type to verify
	// +default="spdxjson"
	predicateType string,
) (string, error) {
	return dag.Container().
		From("gcr.io/projectsigstore/cosign:latest").
		WithMountedSecret("/cosign.pub", publicKey).
		WithExec([]string{
			"cosign", "verify-attestation",
			"--key", "/cosign.pub",
			"--type", predicateType,
			imageRef,
		}).
		Stdout(ctx)
}

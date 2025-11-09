// Dagger module for TruffleHog - secret scanning for any codebase
package main

import (
	"context"
	"dagger/trufflehog/internal/dagger"
	"fmt"
)

type Trufflehog struct{}

// Scan scans a directory for secrets (works with any programming language)
func (m *Trufflehog) Scan(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Output format: json, json-legacy, yaml, text, text-plain, github-actions
	// +default="json"
	format string,
	// Number of concurrent workers
	// +default=10
	concurrency int,
	// Fail on verified secrets
	// +default=true
	failOnVerified bool,
) (string, error) {
	args := []string{
		"trufflehog",
		"filesystem",
		"/src",
		"--json",
		"--no-update",
		fmt.Sprintf("--concurrency=%d", concurrency),
	}

	if !failOnVerified {
		args = append(args, "--no-verification")
	}

	return dag.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

// ScanGit scans a Git repository for secrets (including history)
func (m *Trufflehog) ScanGit(
	ctx context.Context,
	// Git repository URL
	repoUrl string,
	// Branch to scan (optional, scans all if not specified)
	// +optional
	branch string,
	// Maximum depth for git history (0 = scan all history)
	// +default=0
	maxDepth int,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	args := []string{
		"trufflehog",
		"git",
		repoUrl,
		"--json",
		"--no-update",
	}

	if branch != "" {
		args = append(args, "--branch="+branch)
	}

	if maxDepth > 0 {
		args = append(args, "--max-depth="+string(rune(maxDepth+'0')))
	}

	return dag.Container().
		From("trufflesecurity/trufflehog:latest").
		WithExec(args).
		Stdout(ctx)
}

// ScanGithub scans a GitHub repository (requires GITHUB_TOKEN)
func (m *Trufflehog) ScanGithub(
	ctx context.Context,
	// GitHub repository (format: owner/repo)
	repo string,
	// GitHub token for authentication
	token *dagger.Secret,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	return dag.Container().
		From("trufflesecurity/trufflehog:latest").
		WithSecretVariable("GITHUB_TOKEN", token).
		WithExec([]string{
			"trufflehog",
			"github",
			"--repo=" + repo,
			"--json",
			"--no-update",
		}).
		Stdout(ctx)
}

// ScanDocker scans a Docker image for secrets
func (m *Trufflehog) ScanDocker(
	ctx context.Context,
	// Container to scan
	container *dagger.Container,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	tarball := container.AsTarball()

	return dag.Container().
		From("trufflesecurity/trufflehog:latest").
		WithMountedFile("/image.tar", tarball).
		WithExec([]string{
			"trufflehog",
			"docker",
			"--image=file:///image.tar",
			"--json",
			"--no-update",
		}).
		Stdout(ctx)
}

// Verify verifies if detected secrets are valid (makes API calls)
func (m *Trufflehog) Verify(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Only show verified secrets
	// +default=true
	onlyVerified bool,
) (string, error) {
	args := []string{
		"trufflehog",
		"filesystem",
		"/src",
		"--json",
		"--no-update",
	}

	if onlyVerified {
		args = append(args, "--only-verified")
	}

	return dag.Container().
		From("trufflesecurity/trufflehog:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

// Dagger module for .NET SDK operations
// Provides build, test, restore, publish, and format operations
package main

import (
	"context"
	"dagger/dotnet/internal/dagger"
)

type Dotnet struct{}

// Restore restores NuGet packages for a .NET solution or project
func (m *Dotnet) Restore(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Solution or project file to restore
	// +default="."
	project string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (*dagger.Container, error) {
	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", project}), nil
}

// Build builds a .NET solution or project
func (m *Dotnet) Build(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Solution or project file to build
	// +default="."
	project string,
	// Build configuration (Debug or Release)
	// +default="Release"
	configuration string,
	// Additional build arguments
	// +optional
	buildArgs []string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (*dagger.Container, error) {
	args := []string{"dotnet", "build", project, "-c", configuration}
	args = append(args, buildArgs...)

	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", project}).
		WithExec(args), nil
}

// Test runs tests for a .NET project
func (m *Dotnet) Test(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Test project file
	testProject string,
	// Build configuration
	// +default="Release"
	configuration string,
	// Collect code coverage
	// +default=true
	collectCoverage bool,
	// Additional test arguments
	// +optional
	testArgs []string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (string, error) {
	args := []string{"dotnet", "test", testProject, "-c", configuration}

	if collectCoverage {
		args = append(args, "--collect:XPlat Code Coverage", "--results-directory", "/coverage")
	}

	args = append(args, testArgs...)

	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore"}).
		WithExec([]string{"dotnet", "build", "-c", configuration, "--no-restore"}).
		WithExec(args).
		Stdout(ctx)
}

// Publish publishes a .NET project
func (m *Dotnet) Publish(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Project file to publish
	project string,
	// Build configuration
	// +default="Release"
	configuration string,
	// Output directory
	// +default="/app/publish"
	outputDir string,
	// Additional publish arguments
	// +optional
	publishArgs []string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (*dagger.Directory, error) {
	args := []string{"dotnet", "publish", project, "-c", configuration, "-o", outputDir}
	args = append(args, publishArgs...)

	container := dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore"}).
		WithExec([]string{"dotnet", "build", "-c", configuration, "--no-restore"}).
		WithExec(args)

	return container.Directory(outputDir), nil
}

// Format checks or applies code formatting using dotnet format
func (m *Dotnet) Format(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Solution or project file
	// +default="."
	project string,
	// Verify no changes (check mode)
	// +default=true
	verifyNoChanges bool,
	// Verbosity level
	// +default="diagnostic"
	verbosity string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (string, error) {
	args := []string{"dotnet", "format", project}

	if verifyNoChanges {
		args = append(args, "--verify-no-changes")
	}

	args = append(args, "--verbosity", verbosity)

	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", project}).
		WithExec(args).
		Stdout(ctx)
}

// GetCoverage extracts code coverage from test results
func (m *Dotnet) GetCoverage(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Test project file
	testProject string,
	// Build configuration
	// +default="Release"
	configuration string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (string, error) {
	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore"}).
		WithExec([]string{"dotnet", "build", "-c", configuration, "--no-restore"}).
		WithExec([]string{
			"dotnet", "test", testProject,
			"-c", configuration,
			"--no-build",
			"--collect:XPlat Code Coverage",
			"--results-directory", "/coverage",
			"--logger", "trx",
		}).
		WithExec([]string{"sh", "-c", "find /coverage -name 'coverage.cobertura.xml' -exec cat {} \\;"}).
		Stdout(ctx)
}

// BuildWithAnalyzers builds with enhanced security and code analysis
func (m *Dotnet) BuildWithAnalyzers(
	ctx context.Context,
	// Source directory containing .NET project
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Solution or project file
	project string,
	// Build configuration
	// +default="Release"
	configuration string,
	// SDK image version
	// +default="mcr.microsoft.com/dotnet/sdk:8.0"
	sdkImage string,
) (string, error) {
	return dag.Container().
		From(sdkImage).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", project}).
		WithExec([]string{
			"dotnet", "build", project,
			"-c", configuration,
			"/p:TreatWarningsAsErrors=true",
			"/p:EnforceCodeStyleInBuild=true",
			"/p:EnableNETAnalyzers=true",
			"/p:AnalysisLevel=latest",
			"/p:AnalysisMode=AllEnabledByDefault",
		}).
		Stdout(ctx)
}

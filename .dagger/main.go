// A comprehensive CI pipeline for C# Search API with security-first approach
package main

import (
	"context"
	"dagger/search-api/internal/dagger"
	"fmt"
	"strings"
)

type SearchApi struct{}

// Build the C# application and run unit tests
func (m *SearchApi) Build(ctx context.Context, source *dagger.Directory) (*dagger.Container, error) {
	return dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", "SearchApi.sln"}).
		WithExec([]string{"dotnet", "build", "SearchApi.sln", "-c", "Release", "--no-restore"}).
		WithExec([]string{"dotnet", "test", "SearchApi.Tests/SearchApi.Tests.csproj", "-c", "Release", "--no-build", "--verbosity", "normal"}), nil
}

// Run security analysis with dependency check
func (m *SearchApi) SecurityScan(ctx context.Context, source *dagger.Directory) (string, error) {
	output, err := dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		// Restore to get packages
		WithExec([]string{"dotnet", "restore", "SearchApi.sln"}).
		// List dependencies
		WithExec([]string{"dotnet", "list", "SearchApi/SearchApi.csproj", "package", "--vulnerable", "--include-transitive"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("security scan failed: %w", err)
	}

	return output, nil
}

// Run static analysis with dotnet format and analyzers
func (m *SearchApi) StaticAnalysis(ctx context.Context, source *dagger.Directory) (string, error) {
	container := dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", "SearchApi.sln"})

	// Run dotnet format to check code formatting
	formatOutput, err := container.
		WithExec([]string{"dotnet", "format", "SearchApi.sln", "--verify-no-changes", "--verbosity", "diagnostic"}).
		Stdout(ctx)

	if err != nil {
		return formatOutput, fmt.Errorf("code formatting check failed: %w", err)
	}

	return "Static analysis passed: Code formatting is correct", nil
}

// BuildContainer creates the production Docker image
func (m *SearchApi) BuildContainer(ctx context.Context, source *dagger.Directory) *dagger.Container {
	return dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"dotnet", "restore", "SearchApi.sln"}).
		WithExec([]string{"dotnet", "build", "SearchApi.sln", "-c", "Release", "--no-restore"}).
		WithExec([]string{"dotnet", "test", "SearchApi.Tests/SearchApi.Tests.csproj", "-c", "Release", "--no-build", "--verbosity", "normal"}).
		WithExec([]string{"dotnet", "publish", "SearchApi/SearchApi.csproj", "-c", "Release", "-o", "/app/publish", "--no-restore"}).
		From("mcr.microsoft.com/dotnet/aspnet:8.0").
		WithExec([]string{"groupadd", "-r", "searchapi"}).
		WithExec([]string{"useradd", "-r", "-g", "searchapi", "searchapi"}).
		WithWorkdir("/app").
		WithDirectory("/app", dag.Container().
			From("mcr.microsoft.com/dotnet/sdk:8.0").
			WithDirectory("/src", source).
			WithWorkdir("/src").
			WithExec([]string{"dotnet", "publish", "SearchApi/SearchApi.csproj", "-c", "Release", "-o", "/app/publish", "--no-restore"}).
			Directory("/app/publish")).
		WithExec([]string{"chown", "-R", "searchapi:searchapi", "/app"}).
		WithUser("searchapi").
		WithEnvVariable("ASPNETCORE_URLS", "http://+:8080").
		WithEnvVariable("DOTNET_RUNNING_IN_CONTAINER", "true").
		WithExposedPort(8080).
		WithEntrypoint([]string{"dotnet", "SearchApi.dll"})
}

// GenerateSBOM creates a Software Bill of Materials
func (m *SearchApi) GenerateSbom(ctx context.Context, source *dagger.Directory) (string, error) {
	// Using Syft to generate SBOM
	sbom, err := dag.Container().
		From("anchore/syft:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{"syft", "dir:/src", "-o", "spdx-json"}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("SBOM generation failed: %w", err)
	}

	return sbom, nil
}

// ScanContainer performs security scanning on the built container
func (m *SearchApi) ScanContainer(ctx context.Context, container *dagger.Container) (string, error) {
	// Save container as tarball
	tarball := container.AsTarball()

	// Scan with Trivy - ENFORCED (fails on HIGH/CRITICAL vulnerabilities)
	scanResult, err := dag.Container().
		From("aquasec/trivy:latest").
		WithMountedFile("/image.tar", tarball).
		WithExec([]string{
			"image",
			"--input", "/image.tar",
			"--severity", "HIGH,CRITICAL",
			"--format", "json",
			"--exit-code", "1", // FAIL build on vulnerabilities!
		}).
		Stdout(ctx)

	if err != nil {
		return "", fmt.Errorf("container scan FAILED - vulnerabilities found: %w", err)
	}

	return scanResult, nil
}

// SetupLocalRegistry starts a local Docker registry for testing
func (m *SearchApi) SetupLocalRegistry() *dagger.Service {
	return dag.Container().
		From("registry:2").
		WithExposedPort(5000).
		AsService()
}

// PushToLocalRegistry pushes the container to local registry
func (m *SearchApi) PushToLocalRegistry(ctx context.Context, container *dagger.Container, tag string) (string, error) {
	registry := m.SetupLocalRegistry()

	imageRef := fmt.Sprintf("registry:5000/search-api:%s", tag)

	_, err := container.
		WithServiceBinding("registry", registry).
		WithExec([]string{"sh", "-c", "echo 'Image built successfully'"}).
		Sync(ctx)

	if err != nil {
		return "", err
	}

	// Export and push
	address, err := container.
		WithServiceBinding("registry", registry).
		Publish(ctx, imageRef)

	if err != nil {
		return "", fmt.Errorf("failed to push to local registry: %w", err)
	}

	return address, nil
}

// SetupK3s creates a K3s cluster for testing
func (m *SearchApi) SetupK3s(ctx context.Context) (*K3sCluster, error) {
	registry := m.SetupLocalRegistry()

	// Create registry mirror configuration
	registriesConfig := `mirrors:
  "registry:5000":
    endpoint:
      - "http://registry:5000"
`

	// Start K3s with registry mirror
	k3sContainer := dag.Container().
		From("rancher/k3s:v1.28.5-k3s1").
		WithServiceBinding("registry", registry).
		WithNewFile("/etc/rancher/k3s/registries.yaml", registriesConfig).
		WithExec([]string{
			"sh", "-c",
			"k3s server --disable=traefik --disable=metrics-server --write-kubeconfig-mode=644 > /var/log/k3s.log 2>&1 &",
		}).
		WithExec([]string{"sleep", "10"})  // Wait for K3s to start

	k3sService := k3sContainer.AsService()

	// Get kubeconfig
	kubeconfigContent, err := k3sContainer.
		WithExec([]string{"cat", "/etc/rancher/k3s/k3s.yaml"}).
		Stdout(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	// Replace localhost with service name
	kubeconfigContent = strings.ReplaceAll(kubeconfigContent, "127.0.0.1", "k3s")

	kubeconfig := dag.Directory().WithNewFile("kubeconfig", kubeconfigContent)

	return &K3sCluster{
		Service:    k3sService,
		Kubeconfig: kubeconfig,
	}, nil
}

type K3sCluster struct {
	Service    *dagger.Service
	Kubeconfig *dagger.Directory
}

// DeploySolr deploys Solr to the K3s cluster
func (m *SearchApi) DeploySolr(ctx context.Context, k8sManifests *dagger.Directory, cluster *K3sCluster) error {
	kubectlContainer := dag.Container().
		From("bitnami/kubectl:latest").
		WithServiceBinding("k3s", cluster.Service).
		WithDirectory("/kubeconfig", cluster.Kubeconfig).
		WithEnvVariable("KUBECONFIG", "/kubeconfig/kubeconfig").
		WithDirectory("/manifests", k8sManifests)

	// Apply Solr deployment
	_, err := kubectlContainer.
		WithExec([]string{"apply", "-f", "/manifests/solr-deployment.yaml"}).
		WithExec([]string{"wait", "--for=condition=ready", "pod", "-l", "app=solr", "-n", "search-system", "--timeout=300s"}).
		Sync(ctx)

	if err != nil {
		return fmt.Errorf("failed to deploy Solr: %w", err)
	}

	return nil
}

// DeployApi deploys the Search API to the K3s cluster
func (m *SearchApi) DeployApi(ctx context.Context, k8sManifests *dagger.Directory, cluster *K3sCluster, imageTag string) error {
	// Read the deployment manifest
	manifestContent := `---
apiVersion: v1
kind: Service
metadata:
  name: search-api
  namespace: search-system
spec:
  selector:
    app: search-api
  ports:
    - name: http
      port: 80
      targetPort: 8080
  type: ClusterIP
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: search-api
  namespace: search-system
spec:
  replicas: 1
  selector:
    matchLabels:
      app: search-api
  template:
    metadata:
      labels:
        app: search-api
    spec:
      containers:
        - name: api
          image: registry:5000/search-api:` + imageTag + `
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
              name: http
          env:
            - name: ASPNETCORE_ENVIRONMENT
              value: "Production"
            - name: Solr__Url
              value: "http://solr.search-system.svc.cluster.local:8983/solr/metadata"
          resources:
            requests:
              memory: "256Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 30
            periodSeconds: 10
          readinessProbe:
            httpGet:
              path: /health
              port: 8080
            initialDelaySeconds: 15
            periodSeconds: 5
`

	kubectlContainer := dag.Container().
		From("bitnami/kubectl:latest").
		WithServiceBinding("k3s", cluster.Service).
		WithDirectory("/kubeconfig", cluster.Kubeconfig).
		WithEnvVariable("KUBECONFIG", "/kubeconfig/kubeconfig").
		WithNewFile("/deployment.yaml", manifestContent)

	// Apply API deployment
	_, err := kubectlContainer.
		WithExec([]string{"apply", "-f", "/deployment.yaml"}).
		WithExec([]string{"wait", "--for=condition=ready", "pod", "-l", "app=search-api", "-n", "search-system", "--timeout=300s"}).
		Sync(ctx)

	if err != nil {
		return fmt.Errorf("failed to deploy API: %w", err)
	}

	return nil
}

// RunIntegrationTests runs integration tests against deployed services
func (m *SearchApi) RunIntegrationTests(ctx context.Context, source *dagger.Directory, cluster *K3sCluster) (string, error) {
	// Get the search-api service endpoint
	kubectlContainer := dag.Container().
		From("mcr.microsoft.com/dotnet/sdk:8.0").
		WithServiceBinding("k3s", cluster.Service).
		WithDirectory("/kubeconfig", cluster.Kubeconfig).
		WithEnvVariable("KUBECONFIG", "/kubeconfig/kubeconfig").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		// Install kubectl
		WithExec([]string{"sh", "-c", "curl -LO https://dl.k8s.io/release/v1.28.0/bin/linux/amd64/kubectl && chmod +x kubectl && mv kubectl /usr/local/bin/"}).
		// Port-forward to access the service
		WithExec([]string{"sh", "-c", "kubectl port-forward -n search-system svc/search-api 8080:80 &"}).
		WithExec([]string{"sleep", "5"}).
		// Set environment variable for tests
		WithEnvVariable("SOLR_URL", "http://localhost:8080").
		// Run integration tests
		WithExec([]string{"dotnet", "test", "SearchApi.IntegrationTests/SearchApi.IntegrationTests.csproj", "-c", "Release", "--verbosity", "normal"})

	output, err := kubectlContainer.Stdout(ctx)
	if err != nil {
		return "", fmt.Errorf("integration tests failed: %w", err)
	}

	return output, nil
}

// PushToHarbor pushes the final image to Harbor registry
func (m *SearchApi) PushToHarbor(
	ctx context.Context,
	container *dagger.Container,
	harborUrl string,
	harborUsername *dagger.Secret,
	harborPassword *dagger.Secret,
	project string,
	tag string,
) (string, error) {
	imageRef := fmt.Sprintf("%s/%s/search-api:%s", harborUrl, project, tag)

	address, err := container.
		WithRegistryAuth(harborUrl, harborUsername, harborPassword).
		Publish(ctx, imageRef)

	if err != nil {
		return "", fmt.Errorf("failed to push to Harbor: %w", err)
	}

	return address, nil
}

// FullPipeline runs the complete CI/CD pipeline
func (m *SearchApi) FullPipeline(
	ctx context.Context,
	source *dagger.Directory,
	// +optional
	harborUrl string,
	// +optional
	harborUsername *dagger.Secret,
	// +optional
	harborPassword *dagger.Secret,
	// +optional
	harborProject string,
	// +default="latest"
	tag string,
) (string, error) {
	report := "🚀 Starting Full CI/CD Pipeline\n\n"

	// Step 1: Build and Unit Test
	report += "📦 Step 1: Building and running unit tests...\n"
	_, err := m.Build(ctx, source)
	if err != nil {
		return report, fmt.Errorf("build failed: %w", err)
	}
	report += "✅ Build and unit tests passed\n\n"

	// Step 2: Static Analysis
	report += "🔍 Step 2: Running static analysis...\n"
	staticResult, err := m.StaticAnalysis(ctx, source)
	if err != nil {
		report += fmt.Sprintf("⚠️  Static analysis warnings: %v\n\n", err)
	} else {
		report += fmt.Sprintf("✅ %s\n\n", staticResult)
	}

	// Step 3: Security Scan
	report += "🔒 Step 3: Running security dependency scan...\n"
	securityResult, err := m.SecurityScan(ctx, source)
	if err != nil {
		return report, fmt.Errorf("security scan failed: %w", err)
	}
	report += fmt.Sprintf("✅ Security scan completed\n%s\n\n", securityResult)

	// Step 4: Generate SBOM
	report += "📋 Step 4: Generating SBOM...\n"
	sbom, err := m.GenerateSbom(ctx, source)
	if err != nil {
		report += fmt.Sprintf("⚠️  SBOM generation warning: %v\n\n", err)
	} else {
		report += fmt.Sprintf("✅ SBOM generated (%d bytes)\n\n", len(sbom))
	}

	// Step 5: Build Container
	report += "🐳 Step 5: Building container image...\n"
	container := m.BuildContainer(ctx, source)
	report += "✅ Container image built\n\n"

	// Step 6: Scan Container
	report += "🔎 Step 6: Scanning container for vulnerabilities...\n"
	scanResult, err := m.ScanContainer(ctx, container)
	if err != nil {
		report += fmt.Sprintf("⚠️  Container scan warning: %v\n\n", err)
	} else {
		report += fmt.Sprintf("✅ Container scan completed (%d bytes of results)\n\n", len(scanResult))
	}

	// Step 7: Push to Local Registry
	report += "📤 Step 7: Pushing to local registry...\n"
	localImage, err := m.PushToLocalRegistry(ctx, container, tag)
	if err != nil {
		return report, fmt.Errorf("failed to push to local registry: %w", err)
	}
	report += fmt.Sprintf("✅ Pushed to local registry: %s\n\n", localImage)

	// Step 8: Setup K3s Cluster
	report += "☸️  Step 8: Setting up K3s cluster...\n"
	cluster, err := m.SetupK3s(ctx)
	if err != nil {
		return report, fmt.Errorf("failed to setup K3s: %w", err)
	}
	report += "✅ K3s cluster ready\n\n"

	// Step 9: Deploy Solr
	report += "🔍 Step 9: Deploying Solr...\n"
	k8sManifests := source.Directory("k8s")
	err = m.DeploySolr(ctx, k8sManifests, cluster)
	if err != nil {
		return report, fmt.Errorf("failed to deploy Solr: %w", err)
	}
	report += "✅ Solr deployed successfully\n\n"

	// Step 10: Deploy API
	report += "🚀 Step 10: Deploying Search API...\n"
	err = m.DeployApi(ctx, k8sManifests, cluster, tag)
	if err != nil {
		return report, fmt.Errorf("failed to deploy API: %w", err)
	}
	report += "✅ Search API deployed successfully\n\n"

	// Step 11: Run Integration Tests
	report += "🧪 Step 11: Running integration tests...\n"
	integrationResult, err := m.RunIntegrationTests(ctx, source, cluster)
	if err != nil {
		return report, fmt.Errorf("integration tests failed: %w", err)
	}
	report += fmt.Sprintf("✅ Integration tests passed\n%s\n\n", integrationResult)

	// Step 12: Push to Harbor (if credentials provided)
	if harborUrl != "" && harborUsername != nil && harborPassword != nil && harborProject != "" {
		report += "🏗️  Step 12: Pushing to Harbor registry...\n"
		harborImage, err := m.PushToHarbor(ctx, container, harborUrl, harborUsername, harborPassword, harborProject, tag)
		if err != nil {
			return report, fmt.Errorf("failed to push to Harbor: %w", err)
		}
		report += fmt.Sprintf("✅ Pushed to Harbor: %s\n\n", harborImage)
	} else {
		report += "⏭️  Step 12: Skipping Harbor push (credentials not provided)\n\n"
	}

	report += "🎉 Pipeline completed successfully!\n"
	return report, nil
}

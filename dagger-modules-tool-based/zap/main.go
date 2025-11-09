// Dagger module for OWASP ZAP - Dynamic Application Security Testing (DAST)
package main

import (
	"context"
	"dagger/zap/internal/dagger"
)

type Zap struct{}

// BaselineScan runs a ZAP baseline scan against a target (quick passive scan)
func (m *Zap) BaselineScan(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL (e.g., "http://api:8080")
	// +default="http://api:8080"
	targetUrl string,
) (string, error) {
	zapContainer := dag.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithServiceBinding("api", apiService).
		WithMountedCache("/zap/wrk", dag.CacheVolume("zap-reports"))

	_, _ = zapContainer.
		WithExec([]string{
			"zap-baseline.py",
			"-t", targetUrl,
			"-r", "/zap/wrk/report.html",
			"-J", "/zap/wrk/report.json",
			"-w", "/zap/wrk/report.md",
			"-d",
			"-I", // Don't fail on warning
			"-z", "-config api.disablekey=true",
		}).
		Stdout(ctx)

	// Return JSON report
	return zapContainer.
		WithExec([]string{"sh", "-c", "cat /zap/wrk/report.json 2>/dev/null || echo '{}'"}).
		Stdout(ctx)
}

// FullScan runs a full active scan (slower, more comprehensive)
func (m *Zap) FullScan(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
	// Maximum scan duration in minutes
	// +default=10
	maxDuration int,
) (string, error) {
	zapContainer := dag.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithServiceBinding("api", apiService).
		WithMountedCache("/zap/wrk", dag.CacheVolume("zap-reports"))

	_, _ = zapContainer.
		WithExec([]string{
			"zap-full-scan.py",
			"-t", targetUrl,
			"-r", "/zap/wrk/report.html",
			"-J", "/zap/wrk/report.json",
			"-w", "/zap/wrk/report.md",
			"-d",
			"-I",
			"-z", "-config api.disablekey=true",
		}).
		Stdout(ctx)

	return zapContainer.
		WithExec([]string{"sh", "-c", "cat /zap/wrk/report.json 2>/dev/null || echo '{}'"}).
		Stdout(ctx)
}

// ApiScan runs an API-specific scan using OpenAPI/Swagger definition
func (m *Zap) ApiScan(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
	// OpenAPI/Swagger definition file
	apiDefinition *dagger.File,
) (string, error) {
	zapContainer := dag.Container().
		From("ghcr.io/zaproxy/zaproxy:stable").
		WithServiceBinding("api", apiService).
		WithMountedCache("/zap/wrk", dag.CacheVolume("zap-reports")).
		WithMountedFile("/zap/wrk/openapi.json", apiDefinition)

	_, _ = zapContainer.
		WithExec([]string{
			"zap-api-scan.py",
			"-t", "/zap/wrk/openapi.json",
			"-f", "openapi",
			"-r", "/zap/wrk/report.html",
			"-J", "/zap/wrk/report.json",
			"-w", "/zap/wrk/report.md",
			"-d",
			"-I",
			"-z", "-config api.disablekey=true",
		}).
		Stdout(ctx)

	return zapContainer.
		WithExec([]string{"sh", "-c", "cat /zap/wrk/report.json 2>/dev/null || echo '{}'"}).
		Stdout(ctx)
}

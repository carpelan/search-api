// Dagger module for Nuclei - fast, template-based security scanner
package main

import (
	"context"
	"dagger/nuclei/internal/dagger"
)

type Nuclei struct{}

// Scan runs Nuclei with specified templates/tags
func (m *Nuclei) Scan(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL (e.g., "http://api:8080")
	// +default="http://api:8080"
	targetUrl string,
	// Tags to filter templates (e.g., "cve", "owasp", "xss")
	// +default=["owasp"]
	tags []string,
	// Severity levels: info, low, medium, high, critical
	// +default=["high", "critical"]
	severity []string,
) (string, error) {
	args := []string{"nuclei", "-u", targetUrl}

	// Add tags
	if len(tags) > 0 {
		tagStr := ""
		for i, tag := range tags {
			if i > 0 {
				tagStr += ","
			}
			tagStr += tag
		}
		args = append(args, "-tags", tagStr)
	}

	// Add severity
	if len(severity) > 0 {
		sevStr := ""
		for i, sev := range severity {
			if i > 0 {
				sevStr += ","
			}
			sevStr += sev
		}
		args = append(args, "-severity", sevStr)
	}

	args = append(args, "-j", "-silent")

	return dag.Container().
		From("projectdiscovery/nuclei:latest").
		WithServiceBinding("api", apiService).
		WithExec(args).
		Stdout(ctx)
}

// ScanApi runs API-specific security tests
func (m *Nuclei) ScanApi(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
) (string, error) {
	return m.Scan(ctx, apiService, targetUrl, []string{"api", "owasp", "owasp-api-top-10"}, []string{"high", "critical"})
}

// ScanCve scans for known CVEs
func (m *Nuclei) ScanCve(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
) (string, error) {
	return m.Scan(ctx, apiService, targetUrl, []string{"cve"}, []string{"high", "critical"})
}

// ScanWithCustomTemplates scans with custom Nuclei templates
func (m *Nuclei) ScanWithCustomTemplates(
	ctx context.Context,
	// Service to scan
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
	// Directory containing custom .yaml template files
	templates *dagger.Directory,
) (string, error) {
	return dag.Container().
		From("projectdiscovery/nuclei:latest").
		WithServiceBinding("api", apiService).
		WithDirectory("/templates", templates).
		WithExec([]string{
			"nuclei",
			"-u", targetUrl,
			"-t", "/templates",
			"-j",
			"-silent",
		}).
		Stdout(ctx)
}

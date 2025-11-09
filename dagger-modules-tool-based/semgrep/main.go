// Dagger module for Semgrep - multi-language SAST scanner
package main

import (
	"context"
	"dagger/semgrep/internal/dagger"
)

type Semgrep struct{}

// Scan runs Semgrep SAST analysis on source code (works with 30+ languages)
func (m *Semgrep) Scan(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Rule configs (e.g., "p/security-audit", "p/owasp-top-ten", "p/csharp", "p/python")
	// +default=["auto"]
	configs []string,
	// Severity levels to report: INFO, WARNING, ERROR
	// +default=["ERROR", "WARNING"]
	severity []string,
	// Output format: json, sarif, text, gitlab-sast, junit-xml
	// +default="json"
	format string,
	// Exclude patterns (e.g., "*.Tests", "test/", "node_modules/")
	// +optional
	exclude []string,
) (string, error) {
	args := []string{"semgrep"}

	// Add configs
	for _, config := range configs {
		args = append(args, "--config="+config)
	}

	// Add severity levels
	for _, sev := range severity {
		args = append(args, "--severity="+sev)
	}

	// Add excludes
	for _, exc := range exclude {
		args = append(args, "--exclude="+exc)
	}

	// Add format
	if format == "sarif" {
		args = append(args, "--sarif", "--output=/tmp/semgrep-results.sarif")
	} else {
		args = append(args, "--"+format)
	}

	// Disable metrics
	args = append(args, "--metrics=off", ".")

	container := dag.Container().
		From("returntocorp/semgrep:latest").
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec(args)

	if format == "sarif" {
		return container.
			WithExec([]string{"cat", "/tmp/semgrep-results.sarif"}).
			Stdout(ctx)
	}

	return container.Stdout(ctx)
}

// ScanWithCustomRules scans with custom Semgrep rules
func (m *Semgrep) ScanWithCustomRules(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Directory containing custom .yaml/.yml rule files
	rules *dagger.Directory,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	args := []string{
		"semgrep",
		"--config=/rules",
		"--" + format,
		"--metrics=off",
		".",
	}

	return dag.Container().
		From("returntocorp/semgrep:latest").
		WithDirectory("/src", source).
		WithDirectory("/rules", rules).
		WithWorkdir("/src").
		WithExec(args).
		Stdout(ctx)
}

// ScanCi runs Semgrep in CI mode (managed rulesets, comment on PRs)
// Requires SEMGREP_APP_TOKEN for managed scanning
func (m *Semgrep) ScanCi(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Semgrep App token for managed rulesets
	appToken *dagger.Secret,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	return dag.Container().
		From("returntocorp/semgrep:latest").
		WithSecretVariable("SEMGREP_APP_TOKEN", appToken).
		WithDirectory("/src", source).
		WithWorkdir("/src").
		WithExec([]string{
			"semgrep", "ci",
			"--" + format,
		}).
		Stdout(ctx)
}

// ScanLanguage scans with language-specific rulesets
func (m *Semgrep) ScanLanguage(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Language: python, javascript, typescript, java, go, ruby, php, csharp, etc.
	language string,
	// Include security audit rules
	// +default=true
	securityAudit bool,
	// Include OWASP Top 10 rules
	// +default=true
	owaspTopTen bool,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	configs := []string{"p/" + language}

	if securityAudit {
		configs = append(configs, "p/security-audit")
	}

	if owaspTopTen {
		configs = append(configs, "p/owasp-top-ten")
	}

	return m.Scan(ctx, source, configs, []string{"ERROR", "WARNING"}, format, nil)
}

// ScanXss scans specifically for XSS vulnerabilities
func (m *Semgrep) ScanXss(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	return m.Scan(ctx, source, []string{"p/xss"}, []string{"ERROR", "WARNING"}, format, nil)
}

// ScanSqlInjection scans for SQL injection vulnerabilities
func (m *Semgrep) ScanSqlInjection(
	ctx context.Context,
	// Source directory to scan
	// +optional
	// +defaultPath="."
	source *dagger.Directory,
	// Output format
	// +default="json"
	format string,
) (string, error) {
	return m.Scan(ctx, source, []string{"p/sql-injection"}, []string{"ERROR", "WARNING"}, format, nil)
}

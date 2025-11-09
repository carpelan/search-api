// Dagger module for k6 - load and performance testing
package main

import (
	"context"
	"dagger/k6/internal/dagger"
	"fmt"
)

type K6 struct{}

// Run executes a k6 load test with a provided test script
func (m *K6) Run(
	ctx context.Context,
	// Service to test
	apiService *dagger.Service,
	// k6 test script (.js file)
	testScript *dagger.File,
) (string, error) {
	return dag.Container().
		From("grafana/k6:latest").
		WithServiceBinding("api", apiService).
		WithMountedFile("/test.js", testScript).
		WithExec([]string{"k6", "run", "/test.js"}).
		Stdout(ctx)
}

// LoadTest runs a simple load test against an endpoint
func (m *K6) LoadTest(
	ctx context.Context,
	// Service to test
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
	// Endpoint to test
	// +default="/health"
	endpoint string,
	// Number of virtual users
	// +default=10
	vus int,
	// Test duration (e.g., "30s", "2m")
	// +default="30s"
	duration string,
	// P95 response time threshold in milliseconds
	// +default=500
	p95Threshold int,
	// Maximum error rate (0.0-1.0)
	// +default="0.05"
	maxErrorRate string,
) (string, error) {
	testScript := fmt.Sprintf(`
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: %d,
  duration: '%s',
  thresholds: {
    http_req_duration: ['p(95)<%d'],
    http_req_failed: ['rate<%s'],
  },
};

export default function () {
  let response = http.get('%s%s');
  check(response, {
    'status is 200': (r) => r.status === 200,
    'response time < %dms': (r) => r.timings.duration < %d,
  });
  sleep(1);
}
`, vus, duration, p95Threshold, maxErrorRate, targetUrl, endpoint, p95Threshold, p95Threshold)

	return dag.Container().
		From("grafana/k6:latest").
		WithServiceBinding("api", apiService).
		WithNewFile("/test.js", testScript).
		WithExec([]string{"k6", "run", "/test.js"}).
		Stdout(ctx)
}

// StressTest runs a stress test with ramping VUs
func (m *K6) StressTest(
	ctx context.Context,
	// Service to test
	apiService *dagger.Service,
	// Target URL
	// +default="http://api:8080"
	targetUrl string,
	// Endpoint to test
	// +default="/health"
	endpoint string,
	// Maximum virtual users
	// +default=100
	maxVus int,
	// Ramp-up duration
	// +default="1m"
	rampUp string,
	// Plateau duration
	// +default="2m"
	plateau string,
	// Ramp-down duration
	// +default="1m"
	rampDown string,
) (string, error) {
	testScript := fmt.Sprintf(`
import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: '%s', target: %d },  // Ramp up
    { duration: '%s', target: %d },  // Stay at peak
    { duration: '%s', target: 0 },   // Ramp down
  ],
};

export default function () {
  let response = http.get('%s%s');
  check(response, {
    'status is 200': (r) => r.status === 200,
  });
  sleep(1);
}
`, rampUp, maxVus, plateau, maxVus, rampDown, targetUrl, endpoint)

	return dag.Container().
		From("grafana/k6:latest").
		WithServiceBinding("api", apiService).
		WithNewFile("/test.js", testScript).
		WithExec([]string{"k6", "run", "/test.js"}).
		Stdout(ctx)
}

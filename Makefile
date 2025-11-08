.PHONY: help build test run dagger-init dagger-build dagger-full k8s-deploy clean

# Variables
IMAGE_NAME ?= search-api
IMAGE_TAG ?= latest
HARBOR_URL ?= harbor.your-domain.com
HARBOR_PROJECT ?= search-api

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the C# application
	dotnet build SearchApi.sln -c Release

test: ## Run unit tests
	dotnet test SearchApi.Tests/SearchApi.Tests.csproj -c Release --verbosity normal

test-integration: ## Run integration tests (requires Solr)
	dotnet test SearchApi.IntegrationTests/SearchApi.IntegrationTests.csproj -c Release --verbosity normal

run: ## Run the application locally
	cd SearchApi && dotnet run

restore: ## Restore NuGet packages
	dotnet restore SearchApi.sln

clean: ## Clean build artifacts
	dotnet clean SearchApi.sln
	rm -rf */bin */obj

format: ## Format code
	dotnet format SearchApi.sln

format-check: ## Check code formatting
	dotnet format SearchApi.sln --verify-no-changes

security-scan: ## Scan for vulnerable dependencies
	dotnet list SearchApi/SearchApi.csproj package --vulnerable --include-transitive

# All containerization is handled by Dagger
# No need for separate docker-build, docker-compose, etc.
# Use dagger-* targets instead

dagger-init: ## Initialize Dagger module
	cd .dagger && dagger mod init --sdk=go --source=.

dagger-build: ## Build using Dagger
	dagger call build --source=.

dagger-test: ## Run tests using Dagger
	dagger call build --source=.

dagger-security: ## Run security scans using Dagger
	dagger call security-scan --source=.

dagger-static: ## Run static analysis using Dagger
	dagger call static-analysis --source=.

dagger-sbom: ## Generate SBOM using Dagger
	dagger call generate-sbom --source=.

dagger-scan: ## Scan container using Dagger
	dagger call scan-container --container=$$(dagger call build-container --source=.)

dagger-full: ## Run full Dagger pipeline
	dagger call full-pipeline --source=. --tag=$(IMAGE_TAG)

dagger-full-harbor: ## Run full Dagger pipeline with Harbor push
	dagger call full-pipeline \
		--source=. \
		--harbor-url=$(HARBOR_URL) \
		--harbor-username=env:HARBOR_USERNAME \
		--harbor-password=env:HARBOR_PASSWORD \
		--harbor-project=$(HARBOR_PROJECT) \
		--tag=$(IMAGE_TAG)

k8s-deploy-solr: ## Deploy Solr to Kubernetes
	kubectl apply -f k8s/solr-deployment.yaml

k8s-deploy-api: ## Deploy API to Kubernetes
	envsubst < k8s/api-deployment.yaml | kubectl apply -f -

k8s-deploy: k8s-deploy-solr k8s-deploy-api ## Deploy all to Kubernetes

k8s-delete: ## Delete Kubernetes resources
	kubectl delete -f k8s/api-deployment.yaml --ignore-not-found
	kubectl delete -f k8s/solr-deployment.yaml --ignore-not-found

k8s-logs-api: ## View API logs
	kubectl logs -n search-system -l app=search-api -f

k8s-logs-solr: ## View Solr logs
	kubectl logs -n search-system -l app=solr -f

k8s-port-forward: ## Port forward to API
	kubectl port-forward -n search-system svc/search-api 8080:80

index-data: ## Index Riksarkivet data
	./scripts/index-riksarkivet-data.sh

argo-submit: ## Submit Argo Workflow
	argo submit .argo/workflow-template.yaml \
		--parameter image-tag=$(IMAGE_TAG) \
		--watch

argo-deploy: ## Deploy ArgoCD application
	kubectl apply -f .argo/application.yaml

all: clean restore build test ## Run clean, restore, build, and test

ci: format-check security-scan build test ## Run CI pipeline locally

dev: ## Start development environment
	@echo "Starting Solr container for local development..."
	@docker run -d --name search-api-solr -p 8983:8983 solr:9.4 solr-precreate metadata || echo "Solr already running"
	@echo "Starting API with hot reload..."
	cd SearchApi && dotnet watch run

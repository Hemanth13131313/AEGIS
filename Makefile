.PHONY: dev dev-down lint test opa-test helm-lint redteam redteam-dry-run clean help

dev:  ## Start local development environment
	docker compose -f docker-compose.dev.yml up -d

dev-tools:  ## Start dev environment with optional tooling (Kafka UI etc.)
	docker compose -f docker-compose.dev.yml -f docker-compose.tools.yml up -d

dev-down:  ## Stop local dev environment
	docker compose -f docker-compose.dev.yml down

lint:  ## Run all linters
	cd apps/gateway && go vet ./...
	cd apps/policy-engine && go vet ./...
	uv run --project apps/scanner ruff check apps/scanner
	uv run --project apps/rag-monitor ruff check apps/rag-monitor
	uv run --project apps/redteam ruff check apps/redteam

test:  ## Run all unit tests
	cd apps/gateway && go test ./...
	cd apps/policy-engine && go test ./...
	uv run --project apps/scanner pytest apps/scanner
	uv run --project apps/rag-monitor pytest apps/rag-monitor
	uv run --project apps/redteam pytest apps/redteam

opa-test:  ## Run all OPA policy tests
	opa test infra/policies/ apps/policy-engine/rego/ -v

helm-lint:  ## Lint all Helm charts
	helm lint deploy/helm/aegis-gateway
	helm lint deploy/helm/aegis-scanner
	helm lint deploy/helm/aegis-control-plane

redteam:  ## Run red team against local AEGIS (AEGIS_REDTEAM_TARGET required)
	@if [ -z "$$AEGIS_REDTEAM_TARGET" ]; then echo "Set AEGIS_REDTEAM_TARGET"; exit 1; fi
	uv run --project apps/redteam python cli.py run \
		--target $$AEGIS_REDTEAM_TARGET \
		--cases apps/redteam/testcases/registry.json \
		--output table

redteam-dry:  ## Dry-run red team (validate cases, no HTTP requests)
	uv run --project apps/redteam python cli.py run \
		--dry-run \
		--target http://placeholder \
		--cases apps/redteam/testcases/registry.json

redteam-validate:  ## Validate test case registry against schema
	uv run --project apps/redteam python cli.py validate \
		--cases apps/redteam/testcases/registry.json

clean:  ## Remove build artifacts
	rm -rf apps/gateway/bin apps/policy-engine/bin

help:  ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Demo orchestrator for the best-practices module.
#
# The copy-ready files all live in ./template — copy that entire directory
# into a new Go project and nothing has to be edited.
#
# This Makefile only orchestrates the demo against good/ and bad/.

CONFIG       := -c $(CURDIR)/template/.golangci.yml
SBOM_NAME    ?= best-practices
SBOM_VERSION ?= 0.1.0

.PHONY: help \
        demo-good demo-bad \
        sbom-vuln \
        docker-build \
        install-tools clean

# ---------- entry points ----------
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "To copy into a new project: cp -r template/. /path/to/your/project/"

# ---------- demos ----------
demo-good: ## Run the full pipeline against good/ (should be clean)
	go vet ./good/...
	golangci-lint run $(CONFIG) ./good/...
	go test -race ./good/...
	govulncheck ./good/...
	gosec -quiet ./good/...
	nilaway ./good/...

demo-bad: ## Run the full pipeline against bad/ (expected to flag issues)
	-go vet ./bad/...
	-golangci-lint run $(CONFIG) ./bad/...
	-go test -race ./bad/...
	-govulncheck ./bad/...
	-gosec -quiet ./bad/...
	-nilaway ./bad/...

# ---------- security ----------
sbom-vuln: ## Supply chain scan with syft + grype for go.mod + go.sum (both good/ and bad/ use the same dependencies)
	syft . --source-name $(SBOM_NAME) --source-version $(SBOM_VERSION) -o cyclonedx-json=sbom.cdx.json
	grype sbom.cdx.json

# ---------- container ----------
docker-build: ## Smoke-test the template Dockerfile against good/
	@tmp=$$(mktemp -d) && \
	  trap "rm -rf $$tmp" EXIT && \
	  cp template/Dockerfile template/.dockerignore go.mod go.sum $$tmp/ && \
	  cp -r good $$tmp/ && \
	  sed -i.bak 's|\./\.\.\.|./good|' $$tmp/Dockerfile && rm $$tmp/Dockerfile.bak && \
	  docker build -t best-practices-good $$tmp && \
	  docker run --rm best-practices-good

# ---------- utilities ----------
install-tools: ## Install all CLI tools (delegates to template/Makefile)
	$(MAKE) -C template install-tools

clean: ## Remove generated artifacts
	rm -f sbom.cdx.json grype.sarif good/coverage.out good/junit.xml good/bench.txt

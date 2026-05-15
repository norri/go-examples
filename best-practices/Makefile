# Demo orchestrator for the best-practices module.
#
# The copy-ready files all live in ./template — copy that entire directory
# into a new Go project and nothing has to be edited.
#
# This Makefile only orchestrates the demo against good/ and bad/.

CONFIG := -c $(CURDIR)/template/.golangci.yml

.PHONY: help demo-good demo-bad demo-lint demo-vet demo-sec demo-nilaway install-tools clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
	@echo
	@echo "To copy into a new project: cp -r template/. /path/to/your/project/"

demo-good: ## Run lint + tests + vet against good/ (should be clean)
	go vet ./good/...
	golangci-lint run $(CONFIG) ./good/...
	go test -race ./good/...
	gosec -quiet ./good/...
	nilaway ./good/...

demo-bad: ## Run tools against bad/ (expected to flag issues)
	-go vet ./bad/...
	-golangci-lint run $(CONFIG) ./bad/...
	-go test -race ./bad/...
	-gosec -quiet ./bad/...
	-nilaway ./bad/...

SBOM_NAME    ?= best-practices
SBOM_VERSION ?= 0.1.0

sbom-vuln: ## Generate CycloneDX SBOM with syft and scan with grype
	syft . --source-name $(SBOM_NAME) --source-version $(SBOM_VERSION) -o cyclonedx-json=sbom.cdx.json
	grype sbom.cdx.json

install-tools: ## Install all CLI tools (delegates to template/Makefile)
	$(MAKE) -C template install-tools

clean:
	rm -f sbom.cdx.json good/coverage.out good/junit.xml good/bench.txt

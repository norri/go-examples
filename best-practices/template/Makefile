# Copy this Makefile (plus the sibling .golangci.yml and dependabot.yml)
# into a new Go project. All targets operate on ./... — no demo-specific paths.
#
# Bootstrap with `make install-tools`.

SBOM_NAME    ?= best-practices
SBOM_VERSION ?= 0.1.0

.PHONY: help all \
        fmt \
        vet lint nilaway \
        test test-race cover bench \
        vuln sec sbom-vuln \
        mod-tidy mod-verify mod-outdated \
        install-tools clean

# ---------- entry points ----------
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

all: fmt vet lint test-race vuln sec sbom-vuln mod-tidy ## Run the full pipeline

# ---------- formatting ----------
fmt: ## Format code (gofumpt + goimports)
	gofumpt -l -w .
	goimports -w .

# ---------- static analysis ----------
vet: ## go vet
	go vet ./...

lint: ## golangci-lint (umbrella for many linters)
	golangci-lint run ./...

nilaway: ## nilaway — inter-procedural nil-pointer analysis (beta, opt-in)
	nilaway ./...

# ---------- tests ----------
test: ## Unit tests with coverage
	gotestsum --format pkgname --junitfile junit.xml -- -coverprofile=coverage.out ./...

test-race: ## Unit tests with race detector
	gotestsum --format pkgname --junitfile junit.xml -- -race -coverprofile=coverage.out ./...

cover: test ## Open HTML coverage report
	go tool cover -html=coverage.out

bench: ## Run benchmarks (use benchstat to compare runs)
	go test -bench=. -benchmem -count=10 ./... | tee bench.txt

# ---------- security ----------
vuln: ## govulncheck — known CVE scan
	govulncheck ./...

sec: ## gosec — Go-specific security smells
	gosec -quiet ./...

sbom-vuln: ## Supply chain scan with syft + grype
	syft . --source-name $(SBOM_NAME) --source-version $(SBOM_VERSION) -o cyclonedx-json=sbom.cdx.json
	grype sbom.cdx.json

# ---------- module management ----------
mod-tidy: ## Ensure go.mod / go.sum are in sync with imports
	go mod tidy -diff

mod-verify: ## Verify dependencies match checksum database
	go mod verify

mod-outdated: ## List outdated direct dependencies
	go list -u -m -json all | go-mod-outdated -direct -update

# ---------- utilities ----------
install-tools: ## Install all CLI tools used here
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install github.com/psampaz/go-mod-outdated@latest
	go install gotest.tools/gotestsum@latest
	go install go.uber.org/nilaway/cmd/nilaway@latest
	@echo
	@echo "Install separately:"
	@echo "  - syft:       https://github.com/anchore/syft#installation"
	@echo "  - grype:      https://github.com/anchore/grype#installation"
	@echo "  - benchstat:  go install golang.org/x/perf/cmd/benchstat@latest"

clean: ## Remove generated artifacts
	rm -f coverage.out junit.xml bench.txt sbom.cdx.json gosec.sarif grype.sarif

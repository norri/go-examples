# best-practices — recommended tooling for Go projects

A curated, opinionated toolbox covering **linting, security, supply-chain
protection, code quality, bug prevention, performance, and CI**.

## Copy these into your project

Everything copy-ready lives in [`template/`](./template). One command:

```
cp -r best-practices/template/. /path/to/your/project/
cd /path/to/your/project
make install-tools
make all
```

No paths or imports inside those files reference this repo — they work
unchanged in any module.

| File in `template/`      | Goes to                  |
|--------------------------|--------------------------|
| `Makefile`               | Project root             |
| `.golangci.yml`          | Project root             |
| `.pre-commit-config.yaml`| Project root             |
| `dependabot.yml`         | `.github/dependabot.yml` |

The CI workflow at [`.github/workflows/best_practices.yml`](../.github/workflows/best_practices.yml)
in this repo is also reusable — copy it to `.github/workflows/ci.yml` in your project.

## The two demo apps in this repo

| Path     | Purpose |
|----------|---------|
| [`good/`](./good) | Clean implementation — `make all` runs against this. |
| [`bad/`](./bad)   | Same shape, but with **~19 intentionally planted issues** spanning errcheck, gosec, nilerr, bodyclose, noctx, errorlint, ineffassign, unused, predeclared, nilaway, and more. Each issue is numbered in a `// (N) tool: why` comment, so you can diff the two files and see what each tool catches. |

From the `best-practices/` root:

```
make demo-good    # delegates to good/Makefile — runs the full pipeline
make demo-bad     # runs vet, golangci-lint, gosec, nilaway against bad/
```

Sample `make demo-bad` output:

```
14 issues:
* errorlint: 1
* gocritic: 1
* gosec: 7
* nilerr: 1
* noctx: 1
* revive: 1
* rowserrcheck: 1
* unused: 1
```

The rest of this file documents each tool — what it catches, when to use it.

---

## 1. Linters & static analysis

| Tool | What it catches | Notes |
|------|-----------------|-------|
| **go vet** | Printf format mismatches, lock copies, shadowing, unreachable code. | Ships with Go. Always run. |
| **staticcheck** | Hundreds of correctness, simplification, and style checks (SA/ST/QF/S series). The single most valuable analyzer. | Run via golangci-lint. |
| **golangci-lint** | Umbrella runner for ~80 linters with caching and a single config. | The canonical entry point. See `.golangci.yml`. |
| **revive** | Configurable, faster successor to `golint`. Naming, exported-doc, returns. | Configured under `revive` in `.golangci.yml`. |
| **gocritic** | Diagnostics, style and performance suggestions Go vet doesn't make. | Bundled with golangci-lint. |
| **errcheck** | Forgotten error returns. | Catches the #1 Go bug class. |
| **errorlint** | Wrong `%w` usage, `==` on errors instead of `errors.Is`. | Required if you adopt Go 1.13 error wrapping. |
| **bodyclose / rowserrcheck / sqlclosecheck** | Leaked `http.Response.Body`, unchecked `rows.Err()`, unclosed `sql.Rows`. | Common goroutine / connection leaks. |
| **contextcheck / fatcontext / noctx** | Wrong / leaking context usage and HTTP calls without context. | |
| **exhaustive** | Missing `case` arms in enum-like switches. | Pairs well with typed constants. |
| **nilaway** (Uber) | Inter-procedural nil-pointer analysis — tracks nil across function boundaries (still beta). | Run via `make nilaway`; not part of `make all`, not integrated with golangci-lint. |

## 2. Formatting

| Tool | Purpose |
|------|---------|
| **gofmt** | The canonical formatter. Always. |
| **gofumpt** | Stricter superset of gofmt — removes pointless newlines, normalises imports. |
| **goimports** | Adds/removes imports and groups them. Use `-local` to keep your module's imports in a third group. |

Run `make fmt`.

## 3. Security

| Tool | What it does |
|------|--------------|
| **gosec** | Pattern-based scanner for Go security smells (SQL injection via `fmt.Sprintf`, weak crypto, hard-coded creds, unsafe file perms, `crypto/rand` vs `math/rand`). |
| **govulncheck** | Cross-references your imports + call graph against the [Go vulnerability database](https://pkg.go.dev/vuln/). Lower false-positive rate than dep-scanners because it only flags vulns you actually reach. |
| **CodeQL** | GitHub's semantic analysis. Free for public repos; covers taint flows gosec misses. |
| **gitleaks** | Scans diffs/history for secrets. Wired into pre-commit. |

## 4. Supply-chain protection

| Tool / practice | Why |
|-----------------|-----|
| **`go mod verify`** | Confirms your `vendor`/module cache matches `go.sum`. Fails if a dep was tampered with. |
| **`go mod tidy -diff`** | Fails CI if `go.mod`/`go.sum` are out of sync with imports. |
| **`GOFLAGS=-mod=readonly`** | Refuses to silently mutate `go.mod` during builds. |
| **`GOSUMDB=sum.golang.org`** (default) | Validates module hashes against Google's checksum database. Don't disable. |
| **`GOPROXY=https://proxy.golang.org,direct`** | Default proxy provides immutability + cache. For private modules use `GOPRIVATE=…`. |
| **Dependabot / Renovate** | Automated dep PRs. See `dependabot.example.yml`. |
| **Syft + Grype** | Generate an SBOM (`syft . -o cyclonedx-json`) and scan it (`grype sbom:sbom.cdx.json`). Required for SLSA-style supply-chain assurance. |
| **`go-mod-outdated`** | Lists upgradable direct deps in a friendly table. |
| **OpenSSF Scorecard** | Audits your repo against supply-chain best practices (branch protection, pinned actions, etc.). |

## 5. Testing & bug prevention

| Tool | Purpose |
|------|---------|
| **`go test -race`** | Built-in race detector. **Run in CI always.** |
| **`go test -coverprofile`** | Coverage instrumentation. View with `go tool cover -html`. |
| **`go test -fuzz`** | Native fuzzing since Go 1.18 — finds crashes by mutating inputs. |
| **gotestsum** | Wraps `go test` with grouped/colourised output, rerun-on-fail, and JUnit XML for CI dashboards. Used as the default runner in this Makefile. |
| **testify** | `require`/`assert`/`mock` helpers. Use `testifylint` to catch misuse. |
| **mockgen / moq** | Generate interface mocks from go:generate directives. |

## 6. Performance

| Tool | Purpose |
|------|---------|
| **`go test -bench -benchmem`** | Built-in benchmarks; track allocs/op. |
| **benchstat** | Statistical comparison between two benchmark runs. The right way to claim a speed-up. |
| **pprof** (`net/http/pprof`) | CPU / heap / goroutine / block / mutex profiles. |
| **`go tool trace`** | Execution tracer — goroutine scheduling, GC pauses, syscalls. |
| **perflock** (Linux) | Disable frequency scaling for stable benchmarks. |

## 7. Hooks & local DX

* `.pre-commit-config.yaml` runs gofmt, go vet, go-mod-tidy, golangci-lint and
  gitleaks before each commit. Install once with `pre-commit install`.
* `.editorconfig` (add one) — keeps indent / EOL consistent across editors.

## 8. CI

The dedicated workflow at
[`.github/workflows/best_practices.yml`](../.github/workflows/best_practices.yml)
is a comprehensive reference pipeline scoped to this module:

* Parallel jobs: hygiene (gofumpt / goimports / `go mod tidy -diff` / `go mod verify`),
  golangci-lint, tests (gotestsum + `-race -shuffle=on` + JUnit + coverage),
  govulncheck, gosec with SARIF upload, gitleaks, Syft SBOM + Grype scan.
* CodeQL on push and weekly schedule (skipped on PRs as it's expensive).
* Final aggregate `ci` job — point branch protection at it.
* `concurrency` cancels in-flight runs, `permissions: contents: read` at root.
* All third-party actions pinned by commit SHA.

The repo-wide `test_all.yml` runs the basic lint + test matrix across every
Go example — keep an eye on that one too.

Key principles:

* Cache the Go module cache and build cache (`actions/setup-go` does this).
* Always run `-race` and `govulncheck` on PRs.
* Pin third-party actions by SHA in production repos.
* Require `go mod tidy -diff` to pass — prevents drift.

---

## Quick reference — what to add to a new Go project

1. `.golangci.yml` (start from this one, prune linters).
2. `Makefile` with `fmt / lint / test-race / vuln / sec` targets.
3. `.pre-commit-config.yaml`.
4. `dependabot.yml` at `.github/dependabot.yml`.
5. CI workflow that runs all of the above.
6. `SECURITY.md` describing how to report vulnerabilities.

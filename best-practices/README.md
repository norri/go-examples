# best-practices — recommended tooling for Go projects

A curated, opinionated toolbox covering **linting, security, supply-chain
protection, code quality, bug prevention, performance, and CI**.

**TL;DR**

* Drop-in files in [`template/`](./template) — `Makefile`, `.golangci.yml`, `Dockerfile`, `.dockerignore`.
* One reference CI workflow at [`.github/workflows/best-practices-ci.yml`](../.github/workflows/best-practices-ci.yml) and Dependabot config at [`.github/dependabot.yml`](../.github/dependabot.yml).
* Two demo apps — `good/` (clean) vs `bad/` (19 numbered, intentional issues).

```
cp -r best-practices/template/. /path/to/your/project/
cd /path/to/your/project
make install-tools
make all
```

## Copy these into your project

The two files in [`template/`](./template) work unchanged in any module — no
paths or imports reference this repo. The CI workflow and Dependabot config
live under `.github/` in this repo (they can't sit inside `template/` and
still drive this repo's own CI); copy them separately and tweak the paths.

| Source in this repo | Goes to | Purpose |
|---------------------|---------|---------|
| [`template/Makefile`](./template/Makefile) | Project root | `make help` for the full target list. |
| [`template/.golangci.yml`](./template/.golangci.yml) | Project root | Curated golangci-lint v2 config. |
| [`template/Dockerfile`](./template/Dockerfile) | Project root | Multi-stage build → distroless non-root runtime. Adjust the `go build` path if your main package isn't at `./...`. |
| [`template/.dockerignore`](./template/.dockerignore) | Project root | Keep the Docker build context minimal. |
| [`../.github/workflows/best-practices-ci.yml`](../.github/workflows/best-practices-ci.yml) | `.github/workflows/ci.yml` | Reference CI pipeline. Update `working-directory` and the `paths:` filters to match your project layout. |
| [`../.github/dependabot.yml`](../.github/dependabot.yml) | `.github/dependabot.yml` | Weekly Go module + actions updates. Update `directory:` to match your project layout. |

A `SECURITY.md` describing how to report vulnerabilities isn't shipped — add
one yourself.

## The two demo apps in this repo

| Path     | Purpose |
|----------|---------|
| [`good/`](./good) | Clean implementation — `make demo-good` runs the full pipeline against this. |
| [`bad/`](./bad)   | Same shape, but with **19 intentionally planted issues**. Each issue is numbered in a `// (N) tool: why` comment, so you can diff the two files and see what each tool catches. |

From the `best-practices/` root:

```
make demo-good    # runs the full pipeline (vet, lint, race tests, govulncheck, gosec, nilaway) against good/
make demo-bad     # same pipeline against bad/ — failures expected
```

Sample golangci-lint output against `bad/` (14 of the 19 planted issues —
the rest fire via `go vet`, `gosec` standalone, or `nilaway`):

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

All entries below are enabled in [`template/.golangci.yml`](./template/.golangci.yml)
or run via `make`.

| Tool | What it catches |
|------|-----------------|
| [**go vet**](https://pkg.go.dev/cmd/vet) | Printf format mismatches, lock copies, shadowing, unreachable code. Ships with Go. Always run. |
| [**staticcheck**](https://staticcheck.dev/) | Hundreds of correctness, simplification, and style checks (SA/ST/QF/S series). The single most valuable analyzer. Run via golangci-lint. |
| [**golangci-lint**](https://golangci-lint.run/) | Umbrella runner for ~80 linters with caching and a single config. The canonical entry point. See `.golangci.yml`. |
| [**revive**](https://revive.run/) | Configurable, faster successor to `golint`. Naming, exported-doc, returns. Configured under `revive` in `.golangci.yml`. |
| [**gocritic**](https://go-critic.com/) | Diagnostics, style and performance suggestions Go vet doesn't make. Bundled with golangci-lint. |
| [**errcheck**](https://github.com/kisielk/errcheck) | Forgotten error returns — the #1 Go bug class. |
| [**errorlint**](https://github.com/polyfloyd/go-errorlint) | Wrong `%w` usage; `==` on errors instead of `errors.Is`. Required if you wrap errors with `%w`. |
| [**bodyclose**](https://github.com/timakin/bodyclose) / [**rowserrcheck**](https://github.com/jingyugao/rowserrcheck) / [**sqlclosecheck**](https://github.com/ryanrolds/sqlclosecheck) | Leaked `http.Response.Body`, unchecked `rows.Err()`, unclosed `sql.Rows`. Common goroutine / connection leaks. |
| [**contextcheck**](https://github.com/kkHAIKE/contextcheck) / [**fatcontext**](https://github.com/Crocmagnon/fatcontext) / [**noctx**](https://github.com/sonatard/noctx) | Wrong / leaking context usage and HTTP calls without context. |
| [**exhaustive**](https://github.com/nishanths/exhaustive) | Missing `case` arms in enum-like switches. Pairs well with typed constants. |
| [**nilaway**](https://github.com/uber-go/nilaway) (Uber) | Inter-procedural nil-pointer analysis. Beta, opt-in — run via `make nilaway` (separate from `make all`, not bundled in golangci-lint). |

## 2. Formatting

| Tool | What it does |
|------|--------------|
| [**gofmt**](https://pkg.go.dev/cmd/gofmt) | The canonical formatter. Always. |
| [**gofumpt**](https://github.com/mvdan/gofumpt) | Stricter superset of gofmt — removes pointless newlines, normalises imports. |
| [**goimports**](https://pkg.go.dev/golang.org/x/tools/cmd/goimports) | Adds/removes imports and groups them. Use `-local` to keep your module's imports in a third group. |

Run `make fmt`.

## 3. Security

| Tool | What it does |
|------|--------------|
| [**gosec**](https://github.com/securego/gosec) | Pattern-based scanner for Go security smells (SQL injection via `fmt.Sprintf`, weak crypto, hard-coded creds, unsafe file perms, `crypto/rand` vs `math/rand`). |
| [**govulncheck**](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) | Cross-references your imports + call graph against the [Go vulnerability database](https://pkg.go.dev/vuln/). Lower false-positive rate than dep-scanners because it only flags vulns you actually reach. |
| [**CodeQL**](https://codeql.github.com/) | GitHub's semantic analysis. Free for public repos; covers taint flows gosec misses. CI-only. |
| [**gitleaks**](https://github.com/gitleaks/gitleaks) | Scans diffs/history for secrets. CI-only in this setup. |

## 4. Supply-chain protection

| Practice | Why |
|----------|-----|
| [**`go mod verify`**](https://pkg.go.dev/cmd/go#hdr-Verify_dependencies_have_expected_content) | Confirms your `vendor`/module cache matches `go.sum`. Fails if a dep was tampered with. |
| [**`go mod tidy -diff`**](https://pkg.go.dev/cmd/go#hdr-Add_missing_and_remove_unused_modules) | Fails CI if `go.mod`/`go.sum` are out of sync with imports. |
| **`GOFLAGS=-mod=readonly`** | Refuses to silently mutate `go.mod` during builds. *Recommended env var — not set in this repo's CI.* |
| **`GOSUMDB=sum.golang.org`** (default) | Validates module hashes against Google's checksum database. Don't disable. *Go default — no config needed unless you've overridden it.* |
| **`GOPROXY=https://proxy.golang.org,direct`** | Default proxy provides immutability + cache. For private modules use `GOPRIVATE=…`. *Go default — same caveat as above.* |
| [**Dependabot**](https://docs.github.com/en/code-security/dependabot) / [**Renovate**](https://docs.renovatebot.com/) | Automated dep PRs. See [`../.github/dependabot.yml`](../.github/dependabot.yml) — update `directory:` to match your project layout when copying. |
| [**Syft**](https://github.com/anchore/syft) + [**Grype**](https://github.com/anchore/grype) | Generate an SBOM (`syft . -o cyclonedx-json`) and scan it (`grype sbom:sbom.cdx.json`). Useful for SLSA-style supply-chain attestation. |
| [**`go-mod-outdated`**](https://github.com/psampaz/go-mod-outdated) | Lists upgradable direct deps in a friendly table. Available as `make mod-outdated`; not run in CI. |
| [**OpenSSF Scorecard**](https://github.com/ossf/scorecard) | Audits your repo against supply-chain best practices (branch protection, pinned actions, etc.). *Not configured here — recommended for production repos via the [Scorecard Action](https://github.com/ossf/scorecard-action).* |

## 5. Testing & bug prevention

| Tool | What it does |
|------|--------------|
| [**`go test -race`**](https://go.dev/doc/articles/race_detector) | Built-in race detector. **Run in CI always.** |
| [**`go test -coverprofile`**](https://go.dev/blog/cover) | Coverage instrumentation. View with `go tool cover -html`. |
| [**`go test -fuzz`**](https://go.dev/security/fuzz/) | Native fuzzing since Go 1.18 — finds crashes by mutating inputs. *No fuzz targets ship in this example; add `FuzzXxx(f *testing.F)` functions and run `go test -fuzz=Fuzz`.* |
| [**gotestsum**](https://github.com/gotestyourself/gotestsum) | Wraps `go test` with grouped/colourised output, rerun-on-fail, and JUnit XML for CI dashboards. Default runner in the template Makefile. |
| [**testify**](https://github.com/stretchr/testify) | `require`/`assert`/`mock` helpers. Use [`testifylint`](https://github.com/Antonboom/testifylint) (enabled in `.golangci.yml`) to catch misuse. |

## 6. Performance

`good/bench_test.go` ships two benchmarks (`BenchmarkProcess`, `BenchmarkHashSecret`)
demonstrating `b.Loop()` and `b.ReportAllocs()`. Run `make bench` to capture allocs/op
for [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) comparison.

| Tool | What it does |
|------|--------------|
| [**`go test -bench -benchmem`**](https://pkg.go.dev/testing#hdr-Benchmarks) | Built-in benchmarks; track allocs/op. |
| [**benchstat**](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) | Statistical comparison between two benchmark runs. The right way to claim a speed-up. *Not demonstrated here.* |
| [**pprof**](https://pkg.go.dev/net/http/pprof) (`net/http/pprof`) | CPU / heap / goroutine / block / mutex profiles. *Not exposed in this example.* |
| [**`go tool trace`**](https://pkg.go.dev/cmd/trace) | Execution tracer — goroutine scheduling, GC pauses, syscalls. *Not demonstrated here.* |
| [**perflock**](https://github.com/aclements/perflock) (Linux) | Disable frequency scaling for stable benchmarks. *Not demonstrated here.* |

## 7. Container image

[`template/Dockerfile`](./template/Dockerfile) is a multi-stage build that
ends in a [distroless](https://github.com/GoogleContainerTools/distroless)
non-root runtime image. Key choices:

* **Multi-stage** — build deps and Go toolchain stay out of the final image.
* **`CGO_ENABLED=0`** — fully static binary, no glibc / musl dependency.
* **`-trimpath -ldflags="-s -w"`** — reproducible build + smaller binary.
* **BuildKit cache mounts** for `/go/pkg/mod` and `/root/.cache/go-build`.
* **`gcr.io/distroless/static-debian12:nonroot`** — no shell, no package manager, runs as UID 65532. Smaller attack surface than alpine.
* Paired with [`template/.dockerignore`](./template/.dockerignore) to keep the build context lean.

```
docker build -t myapp .
docker run --rm myapp
```

Pin both base images by digest in production (see comment block at the top
of the Dockerfile).

## 8. CI

The dedicated workflow at
[`.github/workflows/best-practices-ci.yml`](../.github/workflows/best-practices-ci.yml)
is a comprehensive reference pipeline scoped to this module:

* Parallel jobs: hygiene (gofumpt / goimports / `go mod tidy -diff` / `go mod verify`),
  golangci-lint, tests (gotestsum + `-race -shuffle=on` + JUnit + coverage),
  govulncheck, gosec with SARIF upload, gitleaks, Syft SBOM + Grype scan.
* CodeQL on push and weekly schedule (skipped on PRs as it's expensive).
* Final aggregate `ci` job — point branch protection at it.
* `concurrency` cancels in-flight runs, `permissions: contents: read` at root.
* All third-party actions pinned by commit SHA.


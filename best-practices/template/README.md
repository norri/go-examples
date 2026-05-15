# Go project template

Copy the contents of this directory into the root of your Go module.
Nothing in these files references a specific module or path — they work
unchanged.

```
cp -r template/. /path/to/your/project/
cd /path/to/your/project
make install-tools
make all
```

## What's here

| File | Where it goes | Purpose |
|------|---------------|---------|
| `Makefile`               | Project root              | `make help` for the full target list. |
| `.golangci.yml`          | Project root              | Curated golangci-lint v2 config. |
| `.pre-commit-config.yaml`| Project root              | gofmt + go vet + go-mod-tidy + golangci-lint + gitleaks. Activate with `pre-commit install`. |
| `dependabot.yml`         | `.github/dependabot.yml`  | Weekly Go module + actions updates. |

See [`../README.md`](../README.md) for the full tool-by-tool rationale.

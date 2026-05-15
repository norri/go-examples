# syntax=docker/dockerfile:1.7
#
# Multi-stage build for a small, static, non-root Go binary.
# References:
#   - https://docs.docker.com/build/building/best-practices/
#   - https://github.com/GoogleContainerTools/distroless
#
# Build:  docker build -t myapp .
# Run:    docker run --rm myapp
#
# For reproducible builds in production, pin base images by digest:
#   FROM golang:1.26-alpine@sha256:<digest> AS build
#   FROM gcr.io/distroless/static-debian12:nonroot@sha256:<digest>

ARG GO_VERSION=1.26

# ---------- build stage ----------
FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src

# Cache module downloads separately from source so code changes don't
# invalidate the dependency layer.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

# CGO_ENABLED=0  → fully static binary, no glibc dependency.
# -trimpath      → strip absolute paths from the binary (reproducibility).
# -ldflags=-s -w → drop the symbol table and DWARF info (smaller binary).
# Change ./... if your main package lives at a specific path (e.g. ./cmd/app).
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -trimpath -ldflags="-s -w" -o /out/app ./...

# ---------- runtime stage ----------
# distroless/static ships only ca-certificates + tzdata — no shell, no package
# manager, smaller attack surface than alpine. `:nonroot` runs as UID 65532.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/app /app

USER nonroot:nonroot
ENTRYPOINT ["/app"]

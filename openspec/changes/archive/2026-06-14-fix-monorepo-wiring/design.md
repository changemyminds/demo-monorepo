## Context

The monorepo wires two Go services and two Python services to live-at-head shared libs via a Go workspace (`go.work`) and a uv workspace (`pyproject.toml`). On a fresh checkout none of the documented `task` entry points work:

- `task run:go APP=service-a` → Go tries `git ls-remote` for `common@v0.0.0` and fails on auth.
- `task test:go` / `task lint:go` → `go test ./...` from the repo root matches no packages.
- `task test:py` → `uv run pytest` fails with `Failed to spawn: pytest`.

All three were reproduced and root-caused empirically (including a minimal-workspace bisection for the Go issue). The service code, shared libs, and tests are correct; only build/environment wiring is broken.

## Goals / Non-Goals

**Goals:**
- Every documented `task` command works on a fresh checkout with no manual fix-ups.
- Shared libs resolve from local source (live-at-head), never from a remote.
- Fixes are minimal, durable, and don't touch application/runtime code.

**Non-Goals:**
- Publishing the shared libs as independently versioned modules.
- Renaming the module path to match the git remote (`example` vs `changemyminds`) — recorded as accepted risk only.
- Reworking the Docker build pipeline beyond the same root-cause fix.

## Decisions

### Decision 1: Pin the Go shared lib with a `replace` directive

Add to `apps/service-a/go.mod` and `apps/service-b/go.mod`:
```
replace github.com/example/demo-monorepo/libs/go/common => ../../libs/go/common
```

**Why:** Bisection proved that the bare `go.work` `use` stanza resolves the local lib *only* when the app has no other external dependencies. Once an app also requires gin, Go performs full module-graph resolution and tries to resolve the non-existent `common@v0.0.0` via VCS — the workspace does not short-circuit it. A path-based `replace` pins the lib to local source unconditionally and is also honored by the `Dockerfile` build (which copies the lib into context), fixing both local and container builds with one change.

**Alternatives considered:**
- *Leave `go.work` as-is and rely on `use`* — rejected: proven not to work with the current toolchain when external deps are present.
- *Add a `replace` in `go.work` instead of per-module* — rejected: Go refuses a versionless workspace-module replace ("replaced at all versions"), and `go.work` is not copied into every build context.
- *Tag a real `libs/go/common/v0.0.0`* — rejected: defeats live-at-head; couples lib releases to app builds.

### Decision 2: Use module-path patterns in the Go Task targets

Change `test:go` to `go test github.com/example/demo-monorepo/...` and `lint:go` to `go vet github.com/example/demo-monorepo/...`.

**Why:** The repo root has no `go.mod`, so a root-relative `./...` resolves against no module and matches nothing. The module-path pattern expands across all workspace modules. Verified to run and pass all three modules.

**Alternatives considered:**
- *`cd` into each module in a loop* — works but more verbose and easy to drift from the `GO_SERVICES` list; the single pattern is simpler.

### Decision 3: Rebuild `.venv` rather than patch shebangs

`task` documentation/flow should ensure `.venv` is (re)created at the current path. Practically: `rm -rf .venv && uv sync --all-packages`.

**Why:** The directory was renamed `demo_monorepo` → `demo-monorepo`; every console-script shebang in `.venv/bin/` still points at the old interpreter path, so the kernel can't exec them (`Failed to spawn`). `uv sync` alone re-links packages but does not rewrite existing wrapper shebangs — a clean recreate is the reliable fix. The venv is a regenerable artifact, so deletion is safe.

**Alternatives considered:**
- *`uv sync --reinstall`* — may rewrite wrappers but is heavier and less obviously correct than a clean recreate.
- *sed-patch the shebangs* — brittle; breaks again on the next path change.

## Risks / Trade-offs

- **`replace` uses relative paths** → only valid within the monorepo layout. Mitigation: documented as intentional for the live-at-head POC; must be removed/versioned before any independent publish (out of scope here).
- **Module path ≠ git remote** (`example` vs `changemyminds`) → a future `go get` or workspace removal would fail. Mitigation: with `replace` + workspace the path is never fetched; recorded as accepted risk in the proposal.
- **venv path-coupling recurs on rename** → any future move of the repo re-breaks `.venv`. Mitigation: `.venv` is gitignored and cheaply rebuilt with `uv sync`; note it in CONTRIBUTING.
- **Go `Dockerfile` not exercised in this change** → believed fixed by the same `replace` but unverified by a container build. Mitigation: validation step builds at least one Go image if feasible; otherwise flagged as follow-up.

## Migration Plan

1. Add `replace` to both Go app `go.mod` files.
2. Update `Taskfile.yml` Go targets.
3. Rebuild `.venv` (`rm -rf .venv && uv sync --all-packages`).
4. Validate: `task test:go`, `task lint:go`, `task test:py`, `task lint:py` all green; run each of the four services and curl `/healthz`.

Rollback: revert the `go.mod`/`Taskfile.yml` edits; `.venv` change is non-versioned and inert.

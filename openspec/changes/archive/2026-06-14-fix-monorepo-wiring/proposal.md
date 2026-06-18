## Why

The scaffolded monorepo POC does not actually run as documented: the headline command `task run:go APP=service-a` fails because Go tries to fetch the local `libs/go/common` lib from GitHub, `task test:go` matches zero packages, and `task test:py` cannot even spawn `pytest`. The service code and tests are correct — only the build/environment wiring is broken — so the POC is non-functional out of the box for any new contributor.

## What Changes

- **Go local-lib resolution**: Add a `replace` directive in each Go app's `go.mod` pointing `github.com/example/demo-monorepo/libs/go/common` at `../../libs/go/common`. The bare `go.work` `use` stanza does not survive module-graph resolution once an app also requires an external dependency (gin): Go then resolves the non-existent `common@v0.0.0` via VCS and fails. The `replace` pins it to local source unconditionally.
- **Taskfile Go targets**: Change `test:go` and `lint:go` from the root-relative `go test ./...` / `go vet ./...` (which match no packages because the repo root has no `go.mod`) to the module-path pattern `github.com/example/demo-monorepo/...`.
- **Python venv rebuild**: Document/automate recreating `.venv` so `uv run pytest` and `uv run uvicorn` work. The project directory was renamed `demo_monorepo` → `demo-monorepo`, leaving every `.venv` console-script shebang pointing at a non-existent interpreter path.
- **Module-path consistency note**: Record that the Go/Python module paths (`example/demo-monorepo`) differ from the git remote (`changemyminds/demo-monorepo`) as a documented, accepted risk for the live-at-head workspace approach.

## Capabilities

### New Capabilities
- `local-dev-workflow`: The contract that a fresh checkout can build, run, and test every service (2 Go + 2 Python) via the documented `task` commands, with the shared libs resolved live-at-head from local source rather than any remote.

### Modified Capabilities
<!-- None: no existing specs in openspec/specs/. -->

## Impact

- **Files**: `apps/service-a/go.mod`, `apps/service-b/go.mod` (add `replace`); `Taskfile.yml` (`test:go`, `lint:go`); `.venv/` (rebuilt); optionally `README.md`/`CONTRIBUTING.md` and Go `Dockerfile`s (same resolution mechanism).
- **Behavior**: `task run:go`, `task test:go`, `task lint:go`, `task test:py`, `task run:py` all succeed; four services respond on `/healthz`, `/version`, `/`.
- **No application/runtime code changes** — service handlers, shared libs, and tests are untouched.
- **Risk**: `replace` directives are local-path based; they must be stripped or versioned if the lib is ever published independently (not a goal for this live-at-head POC).

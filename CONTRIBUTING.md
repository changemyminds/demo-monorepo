# Contributing

This is a learning POC. The workflow mirrors a real monorepo so the patterns transfer.

## Commit convention (required)

We use [Conventional Commits](https://www.conventionalcommits.org/). `release-please`
computes every version bump from commit messages, so the format is load-bearing.

```
<type>(<scope>): <subject>

feat(service-a): add greeting endpoint        # → minor bump for service-a
fix(service-c): correct version reporting      # → patch bump for service-c
chore(libs): tidy go/common helpers            # → no release (lib is live-at-head)
```

- **types**: `feat fix perf refactor docs test build ci chore revert`
- **scopes**: `service-a service-b service-c service-d libs deploy ci release repo`
- `feat!:` or a `BREAKING CHANGE:` footer triggers a major bump.

`commitlint` (see `commitlint.config.js`) and the PR-title check in CI enforce this.

## Local development

Install [Task](https://taskfile.dev), Go, and [uv](https://docs.astral.sh/uv/).

```bash
task test          # all tests (Go + Python)
task run:go APP=service-a
task run:py APP=service-c
task image APP=service-a VERSION=dev
```

Shared libraries are **live at head**: edit `libs/go/common` or `libs/python/common`
and consumers pick up the change immediately — no version bump, no republish.

The Go apps pin the shared lib to local source via a `replace` directive in each
`apps/service-*/go.mod`. This is required: the bare `go.work` `use` stanza stops
resolving the local lib once an app also requires a third-party module, so Go would
otherwise try to fetch the non-existent `common@v0.0.0` over the network.

## Troubleshooting

- **`uv run` fails with `Failed to spawn: pytest` / `uvicorn`** — the `.venv`'s
  console-script shebangs hard-code an absolute interpreter path, so they break if
  the repo is moved or renamed. Rebuild it:

  ```bash
  rm -rf .venv && uv sync --all-packages
  ```

- **Module path vs git remote** — the Go and Python module paths use
  `github.com/example/demo-monorepo`, which does **not** match this repo's git
  remote. This is an accepted trade-off of the live-at-head workspace: the path is
  never fetched (it resolves to local source via the workspace + `replace`). It
  would only matter if someone `go get`s the modules or drops the workspace wiring.

## Branch & PR

1. Branch off `main`.
2. Make focused changes; keep one deployable concern per PR where possible.
3. Ensure `task lint && task test` pass.
4. Open a PR — CI runs lint/test/build for the **affected** units only.

## How a change ships

```
PR merged to main
  → release-please opens/updates a per-app release PR
    → merge it → tag  <service>-vX.Y.Z
      → CI builds & pushes ghcr.io/OWNER/<service>:X.Y.Z
        → auto-deploy to stage
          → manual, reviewer-gated promotion to prd (same image)
```

See `README.md` and `docs`/`ARCHITECTURE.md` for the full rationale.

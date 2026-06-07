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

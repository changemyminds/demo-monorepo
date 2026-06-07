## Why

We need a runnable, teachable monorepo POC that demonstrates the principles in `ARCHITECTURE.md` — organize by deployment unit, share code live-at-head, keep release separate from deploy, and never point a manifest at an image before it exists. Today the repo only contains the architecture doc; there is nothing for a learner to read, run, or copy. A symmetric **2 Go + 2 Python** example makes the patterns concrete without the complexity of the full reference design.

## What Changes

- Add four deployable apps under `apps/`:
  - `service-a`, `service-b` — Go + Gin
  - `service-c`, `service-d` — Python + FastAPI (managed by `uv`)
  - Each exposes `/healthz` and `/version` and consumes its language's shared lib.
- Add two live-at-head shared libraries under `libs/`:
  - `libs/go/common` (wired via `go.work`)
  - `libs/python/common` (wired via a `uv` workspace)
- Add deployment assets under `deploy/`, demonstrating **two patterns side by side**:
  - **Helm** charts for the Go services (`deploy/charts/`)
  - **Kustomize** overlays for the Python services (`deploy/base/` + overlays)
  - Per-environment version selection in `deploy/envs/{stage,prd}/`.
- Add release tooling: `release-please` in **manifest mode** with four packages (libs are NOT packaged), `commitlint`, and conventional-commit enforcement.
- Add GitHub Actions workflows: PR CI (lint/test/build affected), `release-please` on main, tag-triggered build+push to **ghcr.io**, auto-deploy to **stage**, and a manual `workflow_dispatch` **prd** promotion gated by a GitHub Environment reviewer.
- Add repo governance/meta files: root `README.md`, `CONTRIBUTING.md`, `Taskfile.yml`, `.editorconfig`, `.gitignore`, `.gitattributes`, `.github/CODEOWNERS`.
- **Non-goal (documented next step):** `libs/proto/` cross-language contracts and the ArgoCD/GitOps deploy path are described but not implemented in v1.

## Capabilities

### New Capabilities
- `monorepo-structure`: the deploy-unit-based directory layout and the apps/libs/deploy boundary rules.
- `go-services`: two Go + Gin services with health/version endpoints consuming a shared Go lib.
- `python-services`: two Python + FastAPI services managed by uv consuming a shared Python lib.
- `shared-libraries`: live-at-head shared code via `go.work` and a uv workspace, gated by affected-CI rather than versions.
- `release-and-versioning`: release-please manifest mode producing per-app tags from conventional commits.
- `cicd-deployment`: PR CI, tag-triggered build+push, stage auto-deploy and gated prd promotion, with the image-exists-before-manifest rule.

### Modified Capabilities
<!-- None. No existing specs in openspec/specs/. -->

## Impact

- **New files only** — greenfield scaffold; nothing existing is modified or removed.
- **New external dependencies / tooling:** Go toolchain + Gin, Python + uv + FastAPI/uvicorn, Docker, Helm, Kustomize/kubectl, Task, release-please, commitlint, GitHub Actions, ghcr.io.
- **CI/CD:** introduces the `.github/workflows/` pipeline and a GitHub `production` Environment with required reviewers.
- **Learners:** the repo becomes a self-contained reference others can clone, read, and adapt.

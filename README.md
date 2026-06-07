# DevOps Monorepo — Demo / POC

A small, runnable monorepo that demonstrates the principles in
[`ARCHITECTURE.md`](./ARCHITECTURE.md): organize by **deployment unit**, share code
**live at head**, keep **release separate from deploy**, version each app
independently, and **never point a manifest at an image before it exists**.

It ships a symmetric **2 Go + 2 Python** example so each pattern is shown twice.

> Detailed docs are in [`docs/`](./docs). Start here for the map.

## Layout

```
apps/                      deployable units (image + deploy target)
  service-a/  service-b/   Go + Gin
  service-c/  service-d/   Python + FastAPI (uv)
libs/                      live-at-head shared code (no version, not deployed)
  go/common/               wired via go.work
  python/common/           wired via the uv workspace
deploy/                    how & where to deploy (decoupled from app source)
  charts/                  Helm charts        → Go services
  base/ + overlays         Kustomize          → Python services
  envs/{stage,prd}/        which version → which environment
.github/workflows/         ci · release · build+deploy(stage) · deploy(prd)
go.work · pyproject.toml   live-at-head wiring (Go workspace / uv workspace)
release-please-config.json manifest-mode versioning (4 packages)
Taskfile.yml               zero-magic task runner over native tools
```

| Service | Language | Framework | Deploy pattern |
|---------|----------|-----------|----------------|
| service-a | Go | Gin | Helm |
| service-b | Go | Gin | Helm |
| service-c | Python | FastAPI | Kustomize |
| service-d | Python | FastAPI | Kustomize |

Every service exposes `GET /healthz` (200) and `GET /version` (build version), and
uses its language's shared `common` lib.

## Quick start

```bash
# Go (needs Go + go.work)
task run:go APP=service-a      # http://localhost:8080/healthz

# Python (needs uv)
task sync
task run:py APP=service-c      # http://localhost:8000/healthz

task test                     # all tests
task image APP=service-a VERSION=dev
```

## Release → deploy flow

The version is computed **once** by release-please and flows downstream as a derived
image tag — it is never recomputed.

```
conventional commit on main
  → release-please → tag  service-a-v1.2.3        (only place a version is born)
    → build+push   ghcr.io/OWNER/service-a:1.2.3
      → deploy stage  (helm --set image.tag / kustomize edit set image)  [auto]
        → promote prd  (manual workflow_dispatch + Environment reviewer)  [same image]
```

**Hard rule (ARCHITECTURE.md #7):** deploy jobs are `needs: build`, and tags are
injected at deploy time — manifests in git never hard-code a tag before the image
exists. A failed build skips deploy, so the cluster keeps its last good image.

## What you must substitute to run for real

This is a POC; these are intentional placeholders:

- `OWNER` / `ghcr.io/example` → your GitHub org/registry (workflows, Taskfile, deploy values).
- Cluster context & namespaces in `deploy/envs/*` and the deploy workflows.
- GitHub **`production` Environment** with required reviewers (the prd gate).
- CODEOWNERS teams (`@example/*`).

## Deferred (documented next steps)

- `libs/proto/` cross-language contracts with linked-versions (share *contracts*, not code).
- ArgoCD / GitOps pull-based deploy (this POC implements only the push path).

See [`docs/deployment.md`](./docs/deployment.md) for the two-pattern walkthrough and
[`docs/next-steps.md`](./docs/next-steps.md) for the deferred items.

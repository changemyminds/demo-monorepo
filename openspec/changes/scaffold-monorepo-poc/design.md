## Context

`ARCHITECTURE.md` is a normative design doc but the repo has no implementation. This change scaffolds a learning-oriented POC that realizes the doc's principles in the smallest faithful form: a symmetric **2 Go + 2 Python** monorepo. The doc's hard rules carry over verbatim — especially **rule #7: a manifest must never reference an image before that image is built and pushed.**

The exploration settled four forks:
- **Naming:** flat `service-a..d`. Teaching grouping → `service-a`/`service-b` are Go+Gin; `service-c`/`service-d` are Python+FastAPI.
- **Deploy manifests:** dual pattern — Helm for Go, Kustomize for Python — so both tools are demonstrated against the same `deploy/envs` split.
- **Python web framework:** FastAPI + uvicorn.
- **Registry:** ghcr.io.

## Goals / Non-Goals

**Goals:**
- A clone-and-read reference that demonstrates: deploy-unit organization, live-at-head libs, release≠deploy, per-app independent versioning, and the image-before-manifest ordering rule.
- Every service builds, tests, and runs locally with one `task` command.
- One end-to-end path from `git commit` → release-please tag → image → stage → prd.

**Non-Goals:**
- `libs/proto/` cross-language contracts and linked-versions (documented as a next step).
- ArgoCD / GitOps pull-based deploy (the doc's "方式二"); only the push path is implemented.
- Real business logic, auth, persistence, or the webhook-agent service.
- Production-grade cluster/registry credentials; workflows assume ghcr.io + a placeholder cluster.

## Decisions

### D1. Flat service names, language grouped (`a/b`=Go, `c/d`=Python)
Flat names match the doc's `apps/` flat layout and keep language an implementation detail. Grouping a/b vs c/d (rather than interleaving) makes the two language stacks visually obvious in a teaching repo.
*Alternative:* language-prefixed names (`go-service-a`) — rejected as redundant with the README and tag prefixes.

### D2. Live-at-head via native tooling, no lib versions
Go uses `go.work` referencing `libs/go/common`; Python uses a `uv` workspace with `libs/python/common` as a workspace member. Libs are **not** release-please packages. Their safety net is affected-CI: a change under `libs/**` re-tests all consumers.
*Alternative:* version + publish libs internally — rejected; the doc reserves versioning for shipped/externally-pinned artifacts.

### D3. release-please manifest mode, separate PRs, four packages
`release-please-config.json` lists `service-a`(go), `service-b`(go), `service-c`(python), `service-d`(python) with `separate-pull-requests: true`. Each app gets an independent version and a `<service>-vX.Y.Z` tag. Conventional commits are enforced by commitlint so bumps compute correctly.
*Alternative:* single repo version — rejected; couples unrelated deploys.

### D4. Dual deploy pattern, shared envs split
`deploy/charts/<go-app>` (Helm) and `deploy/base/<py-app>` + `deploy/overlays` (Kustomize) both consume `deploy/envs/{stage,prd}` for version selection. This teaches that "how to deploy" (charts/base) is decoupled from "which version where" (envs), regardless of tool.

### D5. Image tag injected after build (rule #7)
`deploy` jobs declare `needs: build`. Helm uses `--set image.tag=$VERSION`; Kustomize uses `kustomize edit set image …:$VERSION`. Manifests in git never hard-code a tag. A failed build means the deploy step never runs, so the cluster keeps its last good image.

### D6. Environments: stage auto, prd manual+gated (promotion model)
Tag push → build → auto-deploy stage (test gate). Prd is a separate `workflow_dispatch` workflow taking a `version` input, bound to a GitHub `production` Environment with required reviewers. Prd is never tied to the stage run, so it can be promoted independently and doesn't hold a run open waiting for approval.

### D7. Task as the build orchestrator
`Taskfile.yml` wraps native tools (`go test`, `uv run`, `docker build`) — zero-magic, no lock-in, matching the doc.

## Risks / Trade-offs

- **Two deploy tools add cognitive load** → mitigated by keeping each service trivial and documenting "why both" in the README; learners can ignore one.
- **ghcr.io + cluster specifics won't run unmodified everywhere** → workflows and envs use clearly-marked placeholders (`ghcr.io/OWNER/...`, namespace names); README lists exactly what to substitute.
- **release-please depends on disciplined commits** → commitlint + PR-title check enforce it; README documents the convention.
- **No proto/ABI sharing** → explicitly deferred; per-language libs only, noted as the natural next exercise.

## Migration Plan

Greenfield — no rollback of existing behavior needed. Suggested build order: structure & tooling → libs → services → local run (Task) → containerization → release-please → CI → deploy (stage) → prd promotion → docs. Each layer is independently verifiable before the next.

## Open Questions

- Target cluster/namespace conventions for the deploy examples (left as documented placeholders for v1).
- Whether to add a `dev` auto-deploy environment later (doc lists it as optional; out of scope for v1).

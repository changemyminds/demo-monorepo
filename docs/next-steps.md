# Deferred next steps

These are intentionally **out of scope for v1** to keep the POC approachable. Each is
a natural follow-on exercise that the current structure already accommodates.

## 1. `libs/proto/` — share contracts, not code

The doc's rule: cross-language sharing = a shared **contract**, not shared source.
Add `libs/proto/` with `.proto` schemas, generate a Go client and a Python client,
and bind them with release-please **linked-versions** (one version for the
`proto / go-client / py-client` bundle).

```
libs/proto/                # schema (source of truth)
  → go-client  (generated) ┐ linked-versions group "proto-bundle"
  → py-client  (generated) ┘
```

Why deferred: adds buf/protoc tooling and a codegen step. The per-language `common`
libs already demonstrate live-at-head; proto demonstrates the *contract* axis.

## 2. ArgoCD / GitOps (the doc's "方式二")

Today we use the **push** path (CI runs `helm`/`kubectl` against the cluster).
The doc's growth path is **pull**: CI only writes the image tag back to
`deploy/envs/**` after build, and ArgoCD syncs the cluster from git.

Migration is cheap because the front half (build+push) and the
`deploy/charts` + `deploy/envs` layout are **already shared**. To switch:

- Add an ArgoCD `Application` (or `ApplicationSet`) per (app × env).
- Replace the deploy job with a `needs: build` writeback (`yq`/`kustomize edit` +
  commit with `[skip ci]`), or use **ArgoCD Image Updater** to avoid writeback
  entirely (it only writes tags that exist in the registry — race-proof by design).
- `release.yml` already sets `paths-ignore: ['deploy/envs/**']` so writebacks don't
  re-trigger releases.

## 3. Optional `dev` auto-deploy environment

The doc lists `dev` as an optional third environment for continuous deployment off
`main`. Add `deploy/envs/dev/` plus a branch-triggered deploy if desired.

## 4. Finer-grained affected detection

CI currently maps changes to affected services (app + its lib edges) via
`dorny/paths-filter`. For larger repos, consider a build graph tool (the doc names
Moon as a "when it hurts" upgrade) rather than path globs.

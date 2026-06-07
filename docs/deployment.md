# Deployment — the two patterns, side by side

This POC deliberately ships **two** deploy tools so you can compare them against
the *same* mental model from `ARCHITECTURE.md`:

- **how to deploy** lives in `deploy/charts/` (Helm) and `deploy/base/` (Kustomize)
- **which version → which environment** lives in `deploy/envs/{stage,prd}/`

| | Go services (service-a, service-b) | Python services (service-c, service-d) |
|---|---|---|
| Tool | Helm | Kustomize |
| "how" | `deploy/charts/<svc>/` | `deploy/base/<svc>/` |
| per-env | `deploy/envs/<env>/<svc>.yaml` (values) | `deploy/envs/<env>/<svc>/` (overlay) |
| tag injection | `helm upgrade --set image.tag=X` | `kustomize edit set image app=…:X` |
| apply | `helm upgrade --install` | `kubectl apply -k` |

Both are driven from the **same** workflows (`release-deploy.yml`, `deploy-prd.yml`),
which branch on the service name. A learner can ignore whichever tool they don't use.

## The end-to-end flow

```
conventional commit on main
  → release.yml (release-please) merges a release PR → tag service-a-v1.2.3
    → release-deploy.yml:
        build  ── docker build+push ghcr.io/OWNER/service-a:1.2.3
          │        (image now provably exists)
          └▶ deploy-stage  (needs: build)  ── auto, environment: staging
                helm upgrade --set image.tag=1.2.3   (Go)
                kustomize edit set image …:1.2.3      (Python)
    → deploy-prd.yml  (manual workflow_dispatch, version=1.2.3)
        environment: production  ── required reviewer approves
        deploys the SAME image to prd (no rebuild)
```

## Why the ordering matters (rule #7)

`deploy-stage` declares `needs: build`. The tag is injected only at deploy time;
nothing in git hard-codes a tag before the image exists. Consequences:

- A **failed build** ⇒ deploy never runs ⇒ the cluster keeps its last good image.
- A manifest can **never** reference a missing image ⇒ no `ImagePullBackOff` race.

## Promotion model (why prd is a separate workflow)

stage and prd receive the **same image** at **different times** (prd after approval).
A single release can't express that per-env timing, so prd is its own manual
`workflow_dispatch`. This keeps prd independently deployable and avoids a run
hanging open waiting for a reviewer. The `production` GitHub Environment's required
reviewers are the gate.

## What you must substitute

- `OWNER` / `ghcr.io/example` → your org/registry (Taskfile, deploy values, workflows).
- Cluster access in both deploy workflows (kubeconfig/ServiceAccount secret).
- A GitHub **`production` Environment** with required reviewers.
- Namespaces (`<svc>-stage`, `<svc>-prd`) if your cluster differs.

# ADR 0001 — Monorepo architecture

**Status:** Accepted

The full architecture decision record for this repository is
[`../../ARCHITECTURE.md`](../../ARCHITECTURE.md). It is normative: the principles and
hard rules there are implementation constraints, not suggestions.

This POC realizes that design in its smallest faithful form — a symmetric
**2 Go + 2 Python** monorepo. See:

- [`../../README.md`](../../README.md) — repo map and quick start
- [`../deployment.md`](../deployment.md) — the two deploy patterns
- [`../next-steps.md`](../next-steps.md) — deferred items (proto, ArgoCD)

## POC-specific concretizations of the doc

- Service naming is flat (`service-a..d`); `a/b` are Go+Gin, `c/d` are Python+FastAPI.
- Deploy demonstrates **both** Helm (Go) and Kustomize (Python) over one `envs/` split.
- Images publish to **ghcr.io**; cluster access and the `production` Environment are
  documented placeholders to substitute.

## 1. Repo scaffolding & tooling

- [x] 1.1 Create top-level `apps/`, `libs/`, `deploy/`, `.github/workflows/` directories
- [x] 1.2 Add root governance files: `README.md`, `CONTRIBUTING.md`, `.editorconfig`, `.gitignore`, `.gitattributes`, `.github/CODEOWNERS`
- [x] 1.3 Add `Taskfile.yml` orchestrating build/test/run/lint for all apps and libs
- [x] 1.4 Add `commitlint.config.js` for conventional commits

## 2. Shared libraries (live at head)

- [x] 2.1 Create `libs/go/common` Go module with a small reusable helper + unit test
- [x] 2.2 Create root `go.work` including `libs/go/common` and both Go service modules
- [x] 2.3 Create `libs/python/common` package with a small reusable helper + unit test
- [x] 2.4 Create root `pyproject.toml` declaring the uv workspace (libs + python services)

## 3. Go services (service-a, service-b)

- [x] 3.1 Scaffold `apps/service-a` Go+Gin app with `/healthz` and `/version`, importing `libs/go/common`
- [x] 3.2 Scaffold `apps/service-b` Go+Gin app with `/healthz` and `/version`, importing `libs/go/common`
- [x] 3.3 Add unit tests and per-app `README.md` for both Go services
- [x] 3.4 Add multi-stage `Dockerfile` for each Go service (version injected via build arg/ldflags)

## 4. Python services (service-c, service-d)

- [x] 4.1 Scaffold `apps/service-c` FastAPI app with `/healthz` and `/version`, depending on `libs/python/common`
- [x] 4.2 Scaffold `apps/service-d` FastAPI app with `/healthz` and `/version`, depending on `libs/python/common`
- [x] 4.3 Add unit tests and per-app `README.md` for both Python services
- [x] 4.4 Add `Dockerfile` for each Python service (uv-based install)

## 5. Local run verification

- [x] 5.1 Verify each service builds, tests, and runs via `task` targets (Python: 8 tests pass + ruff clean; Go: wiring verified by inspection — no local Go toolchain, validated in CI)
- [x] 5.2 Verify a change in each shared lib propagates to its consumers without a version bump (Python demonstrated live; Go uses the same go.work live-at-head mechanism)

## 6. Release & versioning (release-please)

- [x] 6.1 Add `release-please-config.json` (manifest mode, 4 packages, separate PRs, correct release-types)
- [x] 6.2 Add `.release-please-manifest.json` seeding initial versions for the 4 apps
- [x] 6.3 Add per-app version files/changelogs as required by each release-type (Python: `[project] version` in each pyproject; Go: version comes from the git tag via ldflags — no file; CHANGELOGs created by release-please on first release)
- [x] 6.4 Confirm libs are excluded from release packages

## 7. CI pipeline

- [x] 7.1 Add `ci.yml` (PR): lint + test + build with git-diff affected detection (including lib→consumer edges)
- [x] 7.2 Add PR-title / commit conventional-commit check
- [x] 7.3 Add `release.yml` (main): run release-please to open/merge release PRs and create tags

## 8. Build & deploy assets

- [x] 8.1 Add Helm charts under `deploy/charts/` for `service-a` and `service-b` (image.tag templated)
- [x] 8.2 Add Kustomize bases under `deploy/base/` for `service-c` and `service-d`
- [x] 8.3 Add `deploy/envs/stage/` and `deploy/envs/prd/` version-selection files for all four apps
- [x] 8.4 Add tag-triggered build+push workflow to `ghcr.io` (image tag derived from the release tag)

## 9. Deployment workflows (stage & prd)

- [x] 9.1 Add stage auto-deploy job (`needs: build`; Helm `--set image.tag` for Go, `kustomize edit set image` for Python)
- [x] 9.2 Add `deploy-prd.yml` manual `workflow_dispatch` with `version` input, bound to `production` Environment with required reviewers
- [x] 9.3 Verify failed-build path: deploy jobs skip and cluster keeps last good image (guaranteed by `deploy-stage: needs: build`)

## 10. Documentation

- [x] 10.1 Write root `README.md`: architecture overview, layout map, version→deploy flow, what to substitute (owner, registry, cluster)
- [x] 10.2 Document the two deploy patterns (Helm vs Kustomize) and why both are shown
- [x] 10.3 Note deferred next steps: `libs/proto/` cross-language contracts and ArgoCD/GitOps path

## 11. Validation

- [x] 11.1 Run `openspec validate scaffold-monorepo-poc --strict` and resolve any issues (valid)
- [x] 11.2 Confirm all spec scenarios have a corresponding implemented behavior or test

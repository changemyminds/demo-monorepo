## 1. Fix Go local-lib resolution

- [x] 1.1 Add `replace github.com/example/demo-monorepo/libs/go/common => ../../libs/go/common` to `apps/service-a/go.mod`
- [x] 1.2 Add the same `replace` to `apps/service-b/go.mod`
- [x] 1.3 Verify `task run:go APP=service-a` starts and `curl /healthz`, `/version`, `/?name=Ada` all return 200 with the shared-lib message
- [x] 1.4 Verify `task run:go APP=service-b` likewise

## 2. Fix Taskfile Go targets

- [x] 2.1 Change `test:go` command to `go test github.com/example/demo-monorepo/...`
- [x] 2.2 Change `lint:go` command to `go vet github.com/example/demo-monorepo/...`
- [x] 2.3 Verify `task test:go` runs and passes all three Go modules (service-a, service-b, common)
- [x] 2.4 Verify `task lint:go` reports no errors

## 3. Rebuild Python venv

- [x] 3.1 Recreate the venv: `rm -rf .venv && uv sync --all-packages`
- [x] 3.2 Confirm `.venv/bin/pytest` shebang points at an existing interpreter under the current project path
- [x] 3.3 Verify `task test:py` spawns pytest and all 8 Python tests pass
- [x] 3.4 Verify `task run:py APP=service-c` starts and `curl /healthz`, `/?name=Ada` return 200; repeat for `service-d`

## 4. Document accepted risks / follow-ups

- [x] 4.1 Note in `CONTRIBUTING.md` (or README) that `.venv` must be rebuilt with `uv sync` after any repo move/rename
- [x] 4.2 Record the module-path (`example`) vs git-remote (`changemyminds`) mismatch as an accepted live-at-head risk
- [x] 4.3 Built service-a + service-b images. This uncovered a separate Dockerfile bug (`GOFLAGS=-mod=mod` is illegal in workspace mode); fixed both Go Dockerfiles to build the module directly (workspace off, `replace` resolves the lib). Images build and the service-a container serves `/healthz` + shared-lib greeting.

## 5. Final verification

- [x] 5.1 Run `task test` (Go + Python aggregate) and confirm fully green
- [x] 5.2 Run `task lint` and confirm clean
- [x] 5.3 Confirm all four services respond on `/healthz` end-to-end

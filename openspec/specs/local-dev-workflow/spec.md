# local-dev-workflow Specification

## Purpose

A fresh checkout can build, run, and test every service (2 Go + 2 Python) via the documented `task` commands, with the shared libs resolved live-at-head from local source rather than any remote.

## Requirements

### Requirement: Go services build and run from local shared lib

Each Go service SHALL build, run, and resolve the shared `libs/go/common` lib from local source without any network access, even when the service also depends on third-party modules.

#### Scenario: Run a Go service locally

- **WHEN** a developer runs `task run:go APP=service-a` on a fresh checkout with no network access to the module's import host
- **THEN** the service starts and responds 200 on `GET /healthz`, `GET /version`, and `GET /?name=Ada` returns a message containing `Ada` produced by `libs/go/common`

#### Scenario: Shared lib resolves to local source, never remote

- **WHEN** Go resolves the module graph for any Go app that requires both `libs/go/common` and a third-party module (e.g. gin)
- **THEN** `github.com/example/demo-monorepo/libs/go/common` resolves to the local `libs/go/common` directory and Go SHALL NOT attempt a VCS/`git ls-remote` lookup for it

### Requirement: Go test and lint targets cover all modules

The `task test:go` and `task lint:go` targets SHALL execute against every Go module in the workspace and report a non-empty, passing result.

#### Scenario: Run Go tests across the workspace

- **WHEN** a developer runs `task test:go`
- **THEN** tests for `apps/service-a`, `apps/service-b`, and `libs/go/common` all run and pass, and the command does not fail with "directory prefix does not contain modules"

#### Scenario: Run Go vet across the workspace

- **WHEN** a developer runs `task lint:go`
- **THEN** `go vet` runs against all three Go modules and reports no errors

### Requirement: Python services build and test from a healthy venv

The Python virtual environment SHALL allow `task test:py`, `task run:py`, and `task lint:py` to execute their console-script entry points (`pytest`, `uvicorn`, `ruff`) successfully on a fresh checkout.

#### Scenario: Run Python tests

- **WHEN** a developer runs `task test:py`
- **THEN** `pytest` spawns successfully and all tests for `libs/python/common`, `apps/service-c`, and `apps/service-d` pass (no "Failed to spawn" error)

#### Scenario: Run a Python service locally

- **WHEN** a developer runs `task run:py APP=service-c`
- **THEN** `uvicorn` starts and the service responds 200 on `GET /healthz` and `GET /?name=Ada` returns a message containing `Ada` from `libs/python/common`

#### Scenario: venv interpreter path is valid

- **WHEN** the `.venv` is created for the project at its current directory path
- **THEN** every console-script shebang under `.venv/bin/` points at an interpreter path that exists

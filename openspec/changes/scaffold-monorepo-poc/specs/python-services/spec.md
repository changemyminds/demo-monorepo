## ADDED Requirements

### Requirement: Two Python + FastAPI services managed by uv
The repository SHALL provide two Python services, `service-c` and `service-d`, each built with FastAPI, served by uvicorn, and managed by `uv`, located under `apps/`.

#### Scenario: Services build and serve
- **WHEN** a learner runs `apps/service-c` or `apps/service-d` via uv
- **THEN** uv resolves dependencies from the service's `pyproject.toml`
- **AND** the FastAPI app starts under uvicorn on a configurable port

### Requirement: Health and version endpoints
Each Python service SHALL expose a `/healthz` endpoint returning HTTP 200 and a `/version` endpoint returning the build version.

#### Scenario: Health check responds
- **WHEN** a client sends `GET /healthz`
- **THEN** the service responds with HTTP 200

#### Scenario: Version is reported
- **WHEN** a client sends `GET /version`
- **THEN** the service responds with the version string sourced from package metadata

### Requirement: Python services consume the shared Python lib
Each Python service SHALL depend on `libs/python/common` as a uv workspace member resolved live-at-head.

#### Scenario: Shared lib is used at head
- **WHEN** a developer changes a function in `libs/python/common`
- **THEN** both Python services pick up the change without any version bump or republish
- **AND** at least one endpoint's behavior demonstrably uses the shared lib

### Requirement: Containerization
Each Python service SHALL include a `Dockerfile` that installs dependencies with uv and produces a runnable image.

#### Scenario: Image builds
- **WHEN** the service's Dockerfile is built
- **THEN** a container image is produced that runs the service and answers `/healthz`

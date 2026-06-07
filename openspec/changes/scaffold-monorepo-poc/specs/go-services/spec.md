## ADDED Requirements

### Requirement: Two Go + Gin services
The repository SHALL provide two Go services, `service-a` and `service-b`, each built with the Gin web framework and located under `apps/`.

#### Scenario: Services build and serve
- **WHEN** a learner builds and runs `apps/service-a` or `apps/service-b`
- **THEN** the service starts an HTTP server using Gin
- **AND** the server listens on a configurable port (env var with a sensible default)

### Requirement: Health and version endpoints
Each Go service SHALL expose a `/healthz` endpoint returning HTTP 200 and a `/version` endpoint returning the build version.

#### Scenario: Health check responds
- **WHEN** a client sends `GET /healthz`
- **THEN** the service responds with HTTP 200

#### Scenario: Version is reported
- **WHEN** a client sends `GET /version`
- **THEN** the service responds with the version string injected at build time

### Requirement: Go services consume the shared Go lib
Each Go service SHALL import and use `libs/go/common` via live-at-head resolution.

#### Scenario: Shared lib is used at head
- **WHEN** a developer changes a function in `libs/go/common`
- **THEN** both Go services pick up the change without any version bump or republish
- **AND** at least one endpoint's behavior demonstrably uses the shared lib

### Requirement: Containerization
Each Go service SHALL include a `Dockerfile` producing a runnable image.

#### Scenario: Image builds
- **WHEN** the service's Dockerfile is built
- **THEN** a container image is produced that runs the service and answers `/healthz`

## ADDED Requirements

### Requirement: Deploy-unit directory layout
The repository SHALL organize code by deployment unit, with deployable apps under `apps/`, non-deployed shared code under `libs/`, and deployment descriptions under `deploy/`, decoupled from app source.

#### Scenario: Top-level layout exists
- **WHEN** a learner inspects the repository root
- **THEN** `apps/`, `libs/`, `deploy/`, and `.github/workflows/` directories exist
- **AND** `apps/` contains `service-a`, `service-b`, `service-c`, `service-d`

#### Scenario: Apps are deployable, libs are not
- **WHEN** classifying any directory under `apps/` or `libs/`
- **THEN** every `apps/<name>` has a container build and a deploy target
- **AND** no `libs/<...>` has a container build or deploy target

### Requirement: Root tooling and governance files
The repository SHALL place tool configuration at the root (where lint/format/release tools search) and provide standard governance meta files.

#### Scenario: Required root files present
- **WHEN** a learner inspects the repository root
- **THEN** `README.md`, `CONTRIBUTING.md`, `Taskfile.yml`, `.editorconfig`, `.gitignore`, `.gitattributes`, and `commitlint.config.js` exist
- **AND** `.github/CODEOWNERS` exists defining ownership boundaries

### Requirement: Per-app local documentation
Each app SHALL include its own `README.md` describing how to build, test, and run that app locally.

#### Scenario: App README present
- **WHEN** a learner opens any `apps/<service>` directory
- **THEN** a `README.md` exists with local build/test/run instructions for that service

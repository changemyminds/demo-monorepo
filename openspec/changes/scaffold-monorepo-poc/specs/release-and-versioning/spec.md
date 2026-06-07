## ADDED Requirements

### Requirement: release-please manifest mode with per-app packages
The repository SHALL use release-please in manifest mode with four independent packages — `service-a` (go), `service-b` (go), `service-c` (python), `service-d` (python) — using separate pull requests.

#### Scenario: Config declares four packages
- **WHEN** a learner inspects `release-please-config.json` and `.release-please-manifest.json`
- **THEN** the four app paths are declared with their correct `release-type`
- **AND** `separate-pull-requests` is enabled so each app releases independently

#### Scenario: Per-app tags are produced
- **WHEN** release-please opens and merges a release PR for one app
- **THEN** it creates a tag of the form `<service>-vX.Y.Z` for that app only
- **AND** it updates that app's changelog and version file

### Requirement: Conventional commits are enforced
The repository SHALL enforce conventional commits so release-please computes version bumps correctly.

#### Scenario: commitlint configured
- **WHEN** a learner inspects the repository root
- **THEN** `commitlint.config.js` exists configured for conventional commits
- **AND** CI checks PR titles / commits against the convention

### Requirement: Version is computed once and flows downstream
The version SHALL be computed only by release-please; all downstream image tags SHALL be derived from it, not recomputed.

#### Scenario: Tag drives the image tag
- **WHEN** a `<service>-vX.Y.Z` tag is created
- **THEN** the build workflow produces an image tagged `X.Y.Z` from that tag
- **AND** no other step recomputes the version

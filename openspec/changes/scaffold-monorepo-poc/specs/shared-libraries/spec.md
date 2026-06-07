## ADDED Requirements

### Requirement: Live-at-head Go shared library
The repository SHALL provide `libs/go/common` wired into a root `go.work` so that Go apps reference its source directly without publishing a version.

#### Scenario: go.work links the lib
- **WHEN** a learner inspects the root `go.work`
- **THEN** it includes the `libs/go/common` module and both Go service modules
- **AND** the Go services import `libs/go/common` by module path

### Requirement: Live-at-head Python shared library
The repository SHALL provide `libs/python/common` as a member of a root uv workspace so that Python apps reference its source directly without publishing a version.

#### Scenario: uv workspace links the lib
- **WHEN** a learner inspects the root `pyproject.toml`
- **THEN** it declares a uv workspace including `libs/python/common` and both Python service packages
- **AND** the Python services declare `common` as a workspace dependency

### Requirement: Libraries are not versioned or released
Shared libraries SHALL NOT be release-please packages and SHALL NOT carry their own release versions.

#### Scenario: Libs excluded from release config
- **WHEN** a learner inspects `release-please-config.json`
- **THEN** no `libs/**` path appears as a package

### Requirement: Library changes are gated by affected-CI
A change under `libs/**` SHALL trigger re-testing of all consuming apps in CI.

#### Scenario: Lib change re-tests consumers
- **WHEN** CI detects a change under `libs/go/common`
- **THEN** CI runs tests for both Go services
- **WHEN** CI detects a change under `libs/python/common`
- **THEN** CI runs tests for both Python services

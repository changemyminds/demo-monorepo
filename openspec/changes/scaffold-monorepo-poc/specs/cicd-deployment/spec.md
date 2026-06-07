## ADDED Requirements

### Requirement: PR continuous integration with affected detection
The repository SHALL provide a PR CI workflow that lints, tests, and builds only the affected units (apps and their lib dependencies) using git-diff path filtering.

#### Scenario: Only affected units run
- **WHEN** a PR changes one app
- **THEN** CI lints, tests, and builds that app
- **WHEN** a PR changes a shared lib
- **THEN** CI lints, tests, and builds all consumers of that lib

### Requirement: Tag-triggered image build and push to ghcr.io
The repository SHALL build and push a versioned container image to ghcr.io when a release tag is created.

#### Scenario: Image built from tag
- **WHEN** a `<service>-vX.Y.Z` tag is pushed
- **THEN** the workflow builds the service image and pushes `ghcr.io/<owner>/<service>:X.Y.Z`

### Requirement: Image must exist before any manifest references it
A deploy step SHALL run only after the image build+push succeeds, and manifests SHALL NOT hard-code a tag in git before the image exists.

#### Scenario: Deploy depends on build
- **WHEN** the deploy job is defined
- **THEN** it declares `needs: build` so it runs only after a successful push
- **AND** the tag is injected at deploy time (Helm `--set image.tag`, Kustomize `kustomize edit set image`)

#### Scenario: Failed build does not change the cluster
- **WHEN** the build job fails
- **THEN** the deploy job does not run
- **AND** the cluster continues running its previous image

### Requirement: Dual deploy patterns over a shared env split
The repository SHALL demonstrate Helm (for the Go services) and Kustomize (for the Python services), both selecting versions from `deploy/envs/{stage,prd}`.

#### Scenario: Both patterns present
- **WHEN** a learner inspects `deploy/`
- **THEN** Helm charts exist for the Go services and Kustomize bases/overlays exist for the Python services
- **AND** both consume per-environment version selection from `deploy/envs/`

### Requirement: Stage auto-deploy and gated prd promotion
Stage SHALL deploy automatically after a successful build; prd SHALL be a separate manual `workflow_dispatch` promotion gated by a GitHub Environment with required reviewers, deploying the same image.

#### Scenario: Stage deploys automatically
- **WHEN** an image is successfully built from a tag
- **THEN** the workflow deploys that image to the stage environment without manual action

#### Scenario: Prd requires approval and reuses the image
- **WHEN** an operator triggers the prd workflow with a version input
- **THEN** the run pauses for a required reviewer on the `production` Environment
- **AND** upon approval it deploys the same already-built image (no rebuild) to prd

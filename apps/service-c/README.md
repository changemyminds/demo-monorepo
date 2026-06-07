# service-c (Python + FastAPI, uv)

A minimal FastAPI service. Deployed via **Kustomize** (`deploy/base/service-c` + overlays).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness, returns `200 {"status":"ok"}` |
| GET | `/version` | Version from package metadata (managed by release-please) |
| GET | `/?name=Ada` | Greeting via the shared `libs/python/common` lib |

## Run locally

```bash
# from repo root (uv workspace gives live-at-head `common`)
task sync
task run:py APP=service-c        # http://localhost:8000
# or:
uv run --package service-c uvicorn service_c.main:app --reload
```

`PORT` overrides the listen port in the container (default `8000`).

## Test

```bash
task test:py
# or, just this package:
uv run pytest apps/service-c
```

## Build image

```bash
# build context is the REPO ROOT (workspace + shared lib must be present)
task image APP=service-c VERSION=1.2.3
# or:
docker build -f apps/service-c/Dockerfile -t service-c:1.2.3 .
```

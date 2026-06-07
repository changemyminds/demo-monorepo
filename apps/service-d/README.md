# service-d (Python + FastAPI, uv)

A minimal FastAPI service. Deployed via **Kustomize** (`deploy/base/service-d` + overlays).

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/healthz` | Liveness, returns `200 {"status":"ok"}` |
| GET | `/version` | Version from package metadata (managed by release-please) |
| GET | `/?name=Grace` | Greeting via the shared `libs/python/common` lib |

## Run locally

```bash
task sync
task run:py APP=service-d        # http://localhost:8000
# or:
uv run --package service-d uvicorn service_d.main:app --reload
```

`PORT` overrides the listen port in the container (default `8000`).

## Test

```bash
task test:py
# or, just this package:
uv run pytest apps/service-d
```

## Build image

```bash
task image APP=service-d VERSION=1.2.3
# or:
docker build -f apps/service-d/Dockerfile -t service-d:1.2.3 .
```

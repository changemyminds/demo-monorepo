"""service-c — FastAPI demo service.

Mirrors the Go services' contract: /healthz, /version, and a greeting that
exercises the live-at-head shared lib `common`.
"""

from importlib.metadata import PackageNotFoundError, version as pkg_version

from common import greeting
from fastapi import FastAPI

SERVICE_NAME = "service-c"


def _version() -> str:
    try:
        return pkg_version(SERVICE_NAME)
    except PackageNotFoundError:
        return "dev"


app = FastAPI(title=SERVICE_NAME)


@app.get("/healthz")
def healthz() -> dict[str, str]:
    return {"status": "ok"}


@app.get("/version")
def version() -> dict[str, str]:
    return {"service": SERVICE_NAME, "version": _version()}


@app.get("/")
def root(name: str = "world") -> dict[str, str]:
    return {"message": greeting(name)}

# Build from the REPO ROOT so the uv workspace + shared lib are in context:
#   docker build -f apps/service-c/Dockerfile --build-arg VERSION=1.2.3 .
FROM python:3.11-slim
COPY --from=ghcr.io/astral-sh/uv:0.9.18 /uv /uvx /bin/

WORKDIR /app
ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    PORT=8000

# Workspace metadata + every member (uv resolves the whole workspace).
COPY pyproject.toml uv.lock ./
COPY libs/python/common ./libs/python/common
COPY apps/service-c ./apps/service-c
COPY apps/service-d ./apps/service-d

# Install only service-c and its deps (incl. live-at-head `common`), no dev tools.
RUN uv sync --frozen --no-dev --package service-c

EXPOSE 8000
CMD ["sh", "-c", "uv run --no-dev --package service-c uvicorn service_c.main:app --host 0.0.0.0 --port ${PORT}"]

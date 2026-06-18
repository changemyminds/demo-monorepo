# Build from the REPO ROOT so the app module + the shared lib are in context:
#   docker build -f apps/service-a/Dockerfile --build-arg VERSION=1.2.3 .
FROM golang:1.22-alpine AS build
WORKDIR /src

# Allow go to resolve/record checksums during the image build (downloads gin).
# No go.work here on purpose: the app pins the shared lib via a go.mod `replace`
# (../../libs/go/common), so we build the module directly. Workspace mode is off,
# which is what lets GOFLAGS=-mod=mod be used (it is illegal in workspace mode).
ENV GOFLAGS=-mod=mod

# Copy the app module and the shared lib its replace directive points at.
COPY libs/go ./libs/go
COPY apps/service-a ./apps/service-a

ARG VERSION=dev
RUN cd apps/service-a && CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/app .

# Minimal runtime image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/app /app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app"]

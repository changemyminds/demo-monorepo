# Build from the REPO ROOT so go.work + the shared lib are in context:
#   docker build -f apps/service-a/Dockerfile --build-arg VERSION=1.2.3 .
FROM golang:1.22-alpine AS build
WORKDIR /src

# Allow go to resolve/record checksums during the image build.
ENV GOFLAGS=-mod=mod

# Copy the workspace + every module it references (go.work lists all of them).
COPY go.work ./
COPY libs/go ./libs/go
COPY apps/service-a ./apps/service-a
COPY apps/service-b ./apps/service-b

ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
      -ldflags="-s -w -X main.version=${VERSION}" \
      -o /out/app ./apps/service-a

# Minimal runtime image.
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/app /app
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app"]

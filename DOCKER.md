# Docker Notes

## Dockerfile stages

### 1) Builder stage

`FROM golang:1.26.0-alpine3.22 AS builder`

- Purpose: compile the Go binary.
- This stage has Go toolchain and build dependencies.
- Output binary: `/bin/api`.

### 2) Runtime stage (current: development)

`FROM alpine:3.22`

- Purpose: run the compiled binary.
- Includes shell tools, so `docker exec` terminal works during local development.
- Runs as non-root user (`appuser`).

## Is this runtime block builder code?

No. This block is **runtime only**:

```dockerfile
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/api /app/api
RUN adduser -D -u 65532 appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api"]
```

It starts after build is finished and only runs the app.

## Development vs Production

- Development: Alpine runtime (current file) for easier debugging.
- Production: switch runtime stage to distroless for smaller attack surface.

Production runtime example:

```dockerfile
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=builder /bin/api /app/api
EXPOSE 8080
ENTRYPOINT ["/app/api"]
```

## Commands

Build and run:

```bash
docker compose up --build
```

Run in background:

```bash
docker compose up --build -d
```

View logs:

```bash
docker compose logs -f app
```

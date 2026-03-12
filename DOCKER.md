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

## Additional Development Notes

### Development target with Go toolchain

The Dockerfile now also has a `development` target:

```dockerfile
FROM golang:1.26.0-alpine3.22 AS development
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
ENTRYPOINT ["go", "run", "./cmd/api"]
```

Purpose:

- Run the app with `go run` inside the container.
- Make `docker compose watch` useful for Go development.
- Avoid full image rebuild for normal source-code changes.

### Compose development target

In `docker-compose.yml`, `secret_app` can build with:

```yaml
build:
  context: .
  dockerfile: Dockerfile
  target: development
```

Purpose:

- Use the `development` stage instead of the production-style runtime.
- Keep Go toolchain available inside the container.

### Watch command

Start container first:

```bash
docker compose up -d --build secret_app
```

Then start watch mode:

```bash
docker compose watch
```

Purpose:

- Sync code changes into the container.
- Restart the app when source files change.
- Rebuild only when dependency files change.

### Recommended watch behavior

Example watch rules:

```yaml
develop:
  watch:
    - action: sync+restart
      path: ./cmd
      target: /app/cmd
    - action: sync+restart
      path: ./internal
      target: /app/internal
    - action: rebuild
      path: ./go.mod
    - action: rebuild
      path: ./go.sum
```

Meaning:

- `sync+restart`: copy changed files and restart the container process.
- `rebuild`: rebuild image because dependencies or image definition changed.

## Air Hot Reload Notes

### Purpose

Use `Air` inside the development container so Go source changes do not require a full image rebuild.

### Development Docker target with Air

Example:

```dockerfile
FROM golang:1.26.0-alpine3.22 AS development
WORKDIR /app
RUN apk add --no-cache ca-certificates git
RUN go install github.com/air-verse/air@v1.61.7
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
ENTRYPOINT ["air", "-c", ".air.toml"]
```

### Air config file

File:

```toml
.air.toml
```

Purpose:

- Build the app into `./tmp/api`
- Restart automatically when Go files change
- Keep temporary build artifacts inside `tmp/`

### Compose watch behavior with Air

Recommended idea:

```yaml
develop:
  watch:
    - action: sync
      path: ./cmd
      target: /app/cmd
    - action: sync
      path: ./internal
      target: /app/internal
    - action: rebuild
      path: ./go.mod
    - action: rebuild
      path: ./go.sum
```

Meaning:

- `sync`: just copy files into the container
- `Air` inside the container notices the change and restarts the app
- `rebuild` is still needed for dependency changes

### Commands for Air workflow

Build and start:

```bash
docker compose up -d --build secret_app
```

Start watch:

```bash
docker compose watch
```

View logs:

```bash
docker compose logs -f secret_app
```

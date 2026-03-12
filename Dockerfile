# syntax=docker/dockerfile:1
# Development Dockerfile: uses Alpine runtime (shell available) for local debugging/exec.

FROM golang:1.26.0-alpine3.22 AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o /bin/api ./cmd/api

FROM golang:1.26.0-alpine3.22 AS development
WORKDIR /app
RUN apk add --no-cache ca-certificates git
RUN go install github.com/air-verse/air@v1.61.7
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
ENTRYPOINT ["air", "-c", ".air.toml"]

# Development runtime: uses Alpine for local debugging/exec. For production, switch to a distroless image.
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /bin/api /app/api
RUN adduser -D -u 65532 appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api"]

# Production runtime (uncomment and replace this stage when releasing):
# FROM gcr.io/distroless/static-debian12:nonroot
# WORKDIR /app
# COPY --from=builder /bin/api /app/api
# EXPOSE 8080
# ENTRYPOINT ["/app/api"]

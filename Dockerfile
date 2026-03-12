# syntax=docker/dockerfile:1

FROM golang:1.26.0-alpine3.22 AS builder
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o /bin/api ./cmd/api

FROM golang:1.26.0-alpine3.22 AS development
WORKDIR /app
RUN apk add --no-cache ca-certificates git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
EXPOSE 8080
ENTRYPOINT ["go", "run", "./cmd/api"]

FROM alpine:3.22 AS production
WORKDIR /app
COPY --from=builder /bin/api /app/api
RUN adduser -D -u 65532 appuser
USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/api"]

# Production runtime (alternative):
# FROM gcr.io/distroless/static-debian12:nonroot
# WORKDIR /app
# COPY --from=builder /bin/api /app/api
# EXPOSE 8080
# ENTRYPOINT ["/app/api"]

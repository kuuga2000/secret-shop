docker run --rm \
  -v $(pwd):/app \
  -w /app \
  golang:1.26.0-alpine3.22 \
  go test ./...
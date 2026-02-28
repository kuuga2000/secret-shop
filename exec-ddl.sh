#!/bin/bash
set -euo pipefail

set -a
source .env
set +a

docker run --rm --network host -v "$PWD":/app -w /app postgres:16-alpine \
  sh -c "psql \"$DATABASE_URL\" -c \"SELECT id, name, slug FROM products ORDER BY id;\""

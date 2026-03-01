#!/bin/bash
set -euo pipefail

set -a
source .env
set +a

docker run --rm \
  --network host \
  -v "$PWD/migrations:/migrations" \
  migrate/migrate:v4.18.3 \
  -path=/migrations \
  -database "$DATABASE_URL" \
  up

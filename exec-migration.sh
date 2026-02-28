#!/bin/bash

set -a
source .env
set +a

docker run --rm --network host -v "$PWD":/app -w /app postgres:16-alpine \
  sh -c "psql \"$DATABASE_URL\" -f migrations/0001_cart_constraints_and_indexes.sql && \
         psql \"$DATABASE_URL\" -f migrations/0002_order_items_variant_not_null.sql && \
         psql \"$DATABASE_URL\" -f migrations/0003_seed_products_helmets.sql"

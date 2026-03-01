#!/bin/bash
set -euo pipefail

set -a
source .env
set +a

docker run --rm --network host -v "$PWD":/app -w /app postgres:16-alpine \
  sh -c "psql \"$DATABASE_URL\" -c \"SELECT p.id AS product_id, p.name, p.slug, pv.id AS variant_id, pv.sku, pv.price, pv.stock FROM products p LEFT JOIN product_variants pv ON pv.product_id = p.id ORDER BY p.id, pv.id;\""

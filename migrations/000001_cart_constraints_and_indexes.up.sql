-- Cart constraints and query-performance indexes.
-- Safe to run because IF NOT EXISTS is used.

CREATE UNIQUE INDEX IF NOT EXISTS uq_orders_one_draft_per_customer
ON orders (customer_id)
WHERE status = 'draft';

CREATE INDEX IF NOT EXISTS idx_product_variants_product_id
ON product_variants (product_id);

CREATE INDEX IF NOT EXISTS idx_orders_customer_id_status
ON orders (customer_id, status);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id
ON order_items (order_id);

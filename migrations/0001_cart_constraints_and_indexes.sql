-- Cart constraints and query-performance indexes
-- Safe to run multiple times because all indexes use IF NOT EXISTS.

-- 1) One active draft cart per customer
CREATE UNIQUE INDEX IF NOT EXISTS uq_orders_one_draft_per_customer
ON orders (customer_id)
WHERE status = 'draft';

-- 2) Speed up product -> variants lookup
CREATE INDEX IF NOT EXISTS idx_product_variants_product_id
ON product_variants (product_id);

-- 3) Speed up customer cart/order lookups
CREATE INDEX IF NOT EXISTS idx_orders_customer_id_status
ON orders (customer_id, status);

-- 4) Speed up loading items for a cart/order
CREATE INDEX IF NOT EXISTS idx_order_items_order_id
ON order_items (order_id);

-- Note:
-- product_categories(product_id, category_id) already has an index via its
-- composite PRIMARY KEY, so no extra duplicate index is created here.

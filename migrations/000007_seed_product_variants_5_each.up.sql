-- Expand seeded products to five variants each by adding 4 extra variants per product.
-- Existing base SKUs (HELM-00X) remain in place.

INSERT INTO product_variants (product_id, sku, price, stock)
SELECT
  p.id,
  b.base_sku || v.sku_suffix AS sku,
  GREATEST(b.base_price + v.price_delta, 1000) AS price,
  GREATEST(b.base_stock + v.stock_delta, 0) AS stock
FROM (
  VALUES
    ('street-rider-full-face-helmet', 'HELM-001', 129900, 25),
    ('urban-commuter-open-face-helmet', 'HELM-002', 99900, 30),
    ('adventure-trail-modular-helmet', 'HELM-003', 159900, 20),
    ('racing-pro-carbon-helmet', 'HELM-004', 249900, 12),
    ('retro-classic-helmet', 'HELM-005', 119900, 18),
    ('offroad-mx-helmet', 'HELM-006', 139900, 22),
    ('touring-comfort-helmet', 'HELM-007', 169900, 16),
    ('dual-sport-explorer-helmet', 'HELM-008', 149900, 20),
    ('youth-saferide-helmet', 'HELM-009', 89900, 35),
    ('all-weather-shield-helmet', 'HELM-010', 154900, 14)
) AS b(product_slug, base_sku, base_price, base_stock)
JOIN products p ON p.slug = b.product_slug
CROSS JOIN (
  VALUES
    ('-A', -10000, 8),
    ('-B', -5000, 4),
    ('-C', 5000, -3),
    ('-D', 10000, -6)
) AS v(sku_suffix, price_delta, stock_delta)
ON CONFLICT (sku) DO UPDATE
SET
  product_id = EXCLUDED.product_id,
  price = EXCLUDED.price,
  stock = EXCLUDED.stock,
  updated_at = NOW();

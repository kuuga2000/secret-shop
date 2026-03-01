-- Seed one variant per seeded product (10 rows).
INSERT INTO product_variants (product_id, sku, price, stock)
SELECT
  p.id,
  v.sku,
  v.price,
  v.stock
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
) AS v(product_slug, sku, price, stock)
JOIN products p ON p.slug = v.product_slug
ON CONFLICT (sku) DO UPDATE
SET
  product_id = EXCLUDED.product_id,
  price = EXCLUDED.price,
  stock = EXCLUDED.stock,
  updated_at = NOW();

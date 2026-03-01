-- Map each seeded product to one seeded category (10 rows).
INSERT INTO product_categories (product_id, category_id)
SELECT
  p.id,
  c.id
FROM (
  VALUES
    ('street-rider-full-face-helmet', 'full-face'),
    ('urban-commuter-open-face-helmet', 'open-face'),
    ('adventure-trail-modular-helmet', 'modular'),
    ('racing-pro-carbon-helmet', 'racing'),
    ('retro-classic-helmet', 'retro'),
    ('offroad-mx-helmet', 'off-road'),
    ('touring-comfort-helmet', 'touring'),
    ('dual-sport-explorer-helmet', 'dual-sport'),
    ('youth-saferide-helmet', 'youth'),
    ('all-weather-shield-helmet', 'all-weather')
) AS m(product_slug, category_slug)
JOIN products p ON p.slug = m.product_slug
JOIN categories c ON c.slug = m.category_slug
ON CONFLICT DO NOTHING;

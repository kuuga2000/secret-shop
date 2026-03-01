DELETE FROM product_categories pc
USING products p, categories c
WHERE pc.product_id = p.id
  AND pc.category_id = c.id
  AND p.slug IN (
    'street-rider-full-face-helmet',
    'urban-commuter-open-face-helmet',
    'adventure-trail-modular-helmet',
    'racing-pro-carbon-helmet',
    'retro-classic-helmet',
    'offroad-mx-helmet',
    'touring-comfort-helmet',
    'dual-sport-explorer-helmet',
    'youth-saferide-helmet',
    'all-weather-shield-helmet'
  )
  AND c.slug IN (
    'full-face',
    'open-face',
    'modular',
    'racing',
    'retro',
    'off-road',
    'touring',
    'dual-sport',
    'youth',
    'all-weather'
  );

DELETE FROM categories
WHERE slug IN (
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

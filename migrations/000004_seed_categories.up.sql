INSERT INTO categories (name, slug)
VALUES
  ('Full Face', 'full-face'),
  ('Open Face', 'open-face'),
  ('Modular', 'modular'),
  ('Racing', 'racing'),
  ('Retro', 'retro'),
  ('Off Road', 'off-road'),
  ('Touring', 'touring'),
  ('Dual Sport', 'dual-sport'),
  ('Youth', 'youth'),
  ('All Weather', 'all-weather')
ON CONFLICT (slug) DO UPDATE
SET name = EXCLUDED.name;

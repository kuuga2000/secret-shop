-- Seed initial helmet products.
-- Idempotent: unique conflict on slug updates the existing row.

INSERT INTO products (name, description, slug, is_active)
VALUES
  ('Street Rider Full Face Helmet', 'Aerodynamic full-face helmet for daily street riding.', 'street-rider-full-face-helmet', TRUE),
  ('Urban Commuter Open Face Helmet', 'Lightweight open-face helmet designed for city commuting.', 'urban-commuter-open-face-helmet', TRUE),
  ('Adventure Trail Modular Helmet', 'Modular helmet with wide visor and touring comfort.', 'adventure-trail-modular-helmet', TRUE),
  ('Racing Pro Carbon Helmet', 'Premium carbon-shell helmet for high-speed track performance.', 'racing-pro-carbon-helmet', TRUE),
  ('Retro Classic Helmet', 'Vintage-style helmet with modern safety standards.', 'retro-classic-helmet', TRUE),
  ('Offroad MX Helmet', 'Off-road motocross helmet with extended chin guard and peak.', 'offroad-mx-helmet', TRUE),
  ('Touring Comfort Helmet', 'Long-distance touring helmet with noise-reduction padding.', 'touring-comfort-helmet', TRUE),
  ('Dual Sport Explorer Helmet', 'Dual-sport helmet suitable for both on-road and light trail use.', 'dual-sport-explorer-helmet', TRUE),
  ('Youth SafeRide Helmet', 'Smaller-size helmet focused on comfort and safety for young riders.', 'youth-saferide-helmet', TRUE),
  ('All Weather Shield Helmet', 'Helmet with anti-fog visor and ventilation for all-weather rides.', 'all-weather-shield-helmet', TRUE)
ON CONFLICT (slug) DO UPDATE
SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  is_active = EXCLUDED.is_active,
  updated_at = NOW();

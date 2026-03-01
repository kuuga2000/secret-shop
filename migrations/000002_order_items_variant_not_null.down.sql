DO $$
DECLARE
    fk_name TEXT;
BEGIN
    SELECT c.conname
    INTO fk_name
    FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'public'
      AND t.relname = 'order_items'
      AND c.contype = 'f'
      AND pg_get_constraintdef(c.oid) ILIKE '%(product_variant_id)%'
      AND pg_get_constraintdef(c.oid) ILIKE '%REFERENCES product_variants%';

    IF fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE order_items DROP CONSTRAINT %I', fk_name);
    END IF;

    ALTER TABLE order_items
    ALTER COLUMN product_variant_id DROP NOT NULL;

    ALTER TABLE order_items
    ADD CONSTRAINT order_items_product_variant_id_fkey
    FOREIGN KEY (product_variant_id)
    REFERENCES product_variants(id)
    ON DELETE SET NULL;
END $$;

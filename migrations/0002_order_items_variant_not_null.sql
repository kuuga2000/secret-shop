-- Enforce stronger cart/order integrity:
-- order_items.product_variant_id must always be present.
--
-- Safety behavior:
-- 1) Abort migration if existing NULL rows are found.
-- 2) Set column to NOT NULL.
-- 3) Ensure FK delete action is compatible with NOT NULL (RESTRICT/NO ACTION),
--    not SET NULL.

DO $$
DECLARE
    null_count BIGINT;
    fk_name TEXT;
    fk_delete_action "char";
BEGIN
    -- 1) Pre-check: fail fast if data already violates the new constraint.
    SELECT COUNT(*)
    INTO null_count
    FROM order_items
    WHERE product_variant_id IS NULL;

    IF null_count > 0 THEN
        RAISE EXCEPTION
            'Cannot set order_items.product_variant_id NOT NULL. Found % NULL rows.',
            null_count;
    END IF;

    -- 2) Enforce NOT NULL (idempotent).
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'order_items'
          AND column_name = 'product_variant_id'
          AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE order_items
        ALTER COLUMN product_variant_id SET NOT NULL;
    END IF;

    -- 3) If FK action is SET NULL, replace with RESTRICT-compatible FK.
    SELECT c.conname, c.confdeltype
    INTO fk_name, fk_delete_action
    FROM pg_constraint c
    JOIN pg_class t ON t.oid = c.conrelid
    JOIN pg_namespace n ON n.oid = t.relnamespace
    WHERE n.nspname = 'public'
      AND t.relname = 'order_items'
      AND c.contype = 'f'
      AND pg_get_constraintdef(c.oid) ILIKE '%(product_variant_id)%'
      AND pg_get_constraintdef(c.oid) ILIKE '%REFERENCES product_variants%';

    IF fk_name IS NOT NULL AND fk_delete_action = 'n' THEN
        EXECUTE format('ALTER TABLE order_items DROP CONSTRAINT %I', fk_name);

        ALTER TABLE order_items
        ADD CONSTRAINT order_items_product_variant_id_fkey
        FOREIGN KEY (product_variant_id)
        REFERENCES product_variants(id)
        ON DELETE RESTRICT;
    END IF;
END $$;

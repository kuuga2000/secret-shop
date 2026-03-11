DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'customers'
          AND column_name = 'first_name'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'customers'
          AND column_name = 'firstname'
    ) THEN
        ALTER TABLE customers RENAME COLUMN first_name TO firstname;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'customers'
          AND column_name = 'last_name'
    ) AND NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'customers'
          AND column_name = 'lastname'
    ) THEN
        ALTER TABLE customers RENAME COLUMN last_name TO lastname;
    END IF;
END $$;

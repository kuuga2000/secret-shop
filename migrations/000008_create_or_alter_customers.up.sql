-- Create customers or align an existing customers table to the project convention.
-- Naming convention: keep customer name fields as firstname/lastname.

CREATE TABLE IF NOT EXISTS customers (
    id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email VARCHAR(255) NOT NULL,
    password_hash TEXT NOT NULL,
    firstname VARCHAR(100),
    lastname VARCHAR(100),
    phone VARCHAR(32),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    email_verified_at TIMESTAMPTZ,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

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

ALTER TABLE customers
    ADD COLUMN IF NOT EXISTS firstname VARCHAR(100),
    ADD COLUMN IF NOT EXISTS lastname VARCHAR(100),
    ADD COLUMN IF NOT EXISTS phone VARCHAR(32),
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS email_verified_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS last_login_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

ALTER TABLE customers
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN password_hash SET NOT NULL,
    ALTER COLUMN created_at SET DEFAULT NOW(),
    ALTER COLUMN updated_at SET DEFAULT NOW(),
    ALTER COLUMN is_active SET DEFAULT TRUE;

CREATE UNIQUE INDEX IF NOT EXISTS uq_customers_email_lower
ON customers (LOWER(email));

CREATE INDEX IF NOT EXISTS idx_customers_phone
ON customers (phone)
WHERE phone IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_customers_is_active
ON customers (is_active);

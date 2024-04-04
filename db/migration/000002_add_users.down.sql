-- Drop owner and currency unique key
ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "owner_currency_key";

-- Drop foreign key constraint on accounts
ALTER TABLE IF EXISTS accounts DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

-- Drop table users
DROP TABLE IF EXISTS users;
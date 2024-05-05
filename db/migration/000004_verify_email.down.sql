ALTER TABLE users DROP COLUMN IF EXISTS is_email_verified;

DROP TABLE IF EXISTS verify_emails CASCADE;


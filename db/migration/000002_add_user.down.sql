
ALTER TABLE IF EXISTS "customer_accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";

ALTER TABLE IF EXISTS "customer_accounts" DROP CONSTRAINT IF EXISTS "customer_accounts_customer_name_fkey";

drop table users cascade;
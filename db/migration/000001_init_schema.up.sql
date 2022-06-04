
CREATE TABLE "customer_accounts" (
                                     "id" bigserial PRIMARY KEY,
                                     "customer_name" varchar NOT NULL,
                                     "balance" bigint NOT NULL,
                                     "currency" varchar NOT NULL,
                                     "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "entries" (
                           "id" bigserial PRIMARY KEY,
                           "customer_id" bigint NOT NULL,
                           "amount" bigint NOT NULL,
                           "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE "transaction" (
                               "id" bigserial PRIMARY KEY,
                               "from_customer_accounts" bigint NOT NULL,
                               "to_customer_accounts" bigint NOT NULL,
                               "amount" bigint NOT NULL,
                               "created_at" timestamp NOT NULL DEFAULT 'now()'
);

CREATE INDEX ON "customer_accounts" ("customer_name");

CREATE INDEX ON "entries" ("customer_id");

CREATE INDEX ON "transaction" ("from_customer_accounts");

CREATE INDEX ON "transaction" ("to_customer_accounts");

CREATE INDEX ON "transaction" ("from_customer_accounts", "to_customer_accounts");

-- ALTER TABLE "entries" ADD FOREIGN KEY ("customer_id") REFERENCES "customer_accounts" ("id");
--
-- ALTER TABLE "transaction" ADD FOREIGN KEY ("from_customer_accounts") REFERENCES "customer_accounts" ("id");
--
-- ALTER TABLE "transaction" ADD FOREIGN KEY ("to_customer_accounts") REFERENCES "customer_accounts" ("id");

-- users 
CREATE TABLE users (
    "id" bigserial PRIMARY KEY,
    "email" varchar(255) NOT NULL,
    "amount" int NOT NULL DEFAULT 0,
    "login_token" varchar(512) NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "updated_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "deleted_at" timestamptz
);

-- api keys 
CREATE TABLE api_keys (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "api_key" varchar(512) NOT NULL,
    "spend" int NOT NULL DEFAULT 0,
    "status" smallint NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "updated_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "deleted_at" timestamptz
);

-- account incomes
CREATE TABLE account_incomes (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "recharge_type" smallint  NOT NULL,
    "recharge_value" int  NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "updated_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "deleted_at" timestamptz
);

-- account incomes
CREATE TABLE account_outcomes (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "api_key" bigint  NOT NULL,
    "tokens" int  NOT NULL,
    "fee_rate" int  NOT NULL,
    "cost" int  NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "updated_at" timestamptz NOT NULL DEFAULT TIME('now'),
    "deleted_at" timestamptz
);

-- users 
CREATE TABLE IF NOT EXISTS users (
    "id" bigserial PRIMARY KEY,
    "email" varchar(255) NOT NULL,
    "amount" integer NOT NULL DEFAULT 0,
    "verification_code" varchar(8) NOT NULL,
    "verification_at" timestamptz NOT NULL,
    "session_id" varchar(512), 
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

-- api keys 
CREATE TABLE IF NOT EXISTS api_keys (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "name" varchar(255) NOT NULL,
    "key" varchar(512) NOT NULL,
    "spend" integer NOT NULL DEFAULT 0,
    "status" smallint NOT NULL DEFAULT 0,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

-- account incomes
CREATE TABLE IF NOT EXISTS incomes (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "recharge_type" smallint  NOT NULL,
    "recharge_value" integer  NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

-- account outcomes
CREATE TABLE IF NOT EXISTS outcomes (
    "id" bigserial PRIMARY KEY,
    "user_id" bigint NOT NULL,
    "key_id" bigint NOT NULL,
    "prompt_tokens" integer NOT NULL,
    "completion_tokens" integer NOT NULL,
    "total_tokens" integer NOT NULL,
    "fee_rate" integer NOT NULL,
    "cost" integer NOT NULL,
    "rt" integer NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT now(),
    "updated_at" timestamptz NOT NULL DEFAULT now(),
    "deleted_at" timestamptz
);

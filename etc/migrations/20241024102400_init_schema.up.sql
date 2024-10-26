CREATE TABLE IF NOT EXISTS "user" (
    "id"                    SERIAL PRIMARY KEY,
    "email"                 VARCHAR(255) NOT NULL,
    "password"              VARCHAR(255) NOT NULL, 
    "first_name"            VARCHAR(255) NOT NULL DEFAULT '',
    "last_name"             VARCHAR(255) NOT NULL DEFAULT '',
    "gender"                INTEGER NOT NULL DEFAULT 0, -- 0=male 1=female
    "status"                SMALLINT NOT NULL DEFAULT 0, -- 0=active 1=inactive 2=suspens
    "deleted"               SMALLINT NOT NULL DEFAULT 0, -- 0=normal 1=deleted
    "updated_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "key_user_email" ON "user" ("email");

CREATE TABLE IF NOT EXISTS "wallet" (
    "id"                    SERIAL PRIMARY KEY,
    "user_id"               INTEGER NOT NULL,
    "balance"               DECIMAL(12, 2) NOT NULL default 0,
    "deleted"               SMALLINT NOT NULL DEFAULT 0,
    "updated_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX "key_wallet_user_id" ON "wallet" ("user_id");

CREATE TABLE IF NOT EXISTS "wallet_record" (
    "id"                    SERIAL PRIMARY KEY,
    "wallet_id"             INTEGER NOT NULL,
    "amount"                DECIMAL(12, 2) NOT NULL DEFAULT 0,
    "transaction_type"      SMALLINT NOT NULL DEFAULT 0, 
    "describe"              TEXT,
    "deleted"               SMALLINT NOT NULL DEFAULT 0,
    "updated_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "created_at"            TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "idx_wallet_record_wallet_id" ON "wallet_record" ("wallet_id");

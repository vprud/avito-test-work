-- \c avito_processing;
DROP TYPE IF EXISTS trans_type;
DROP TABLE IF EXISTS account;
DROP TABLE IF EXISTS fct_transcation;
CREATE TYPE trans_type AS ENUM ('accrual', 'redeem');
CREATE TABLE account (
	id BIGSERIAL PRIMARY KEY,
	balance NUMERIC (16, 3) NOT NULL DEFAULT 0.000 CHECK (balance >= 0.000),
    created_dt TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE fct_transcation (
	id SERIAL PRIMARY KEY,
    trans_dt TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    account_id BIGINT REFERENCES account ON DELETE CASCADE,
    doc_num BIGINT DEFAULT -999, -- redeem_id
    type trans_type,
	amount NUMERIC(16, 3) NOT NULL
);
INSERT INTO account VALUES(-999, DEFAULT, DEFAULT);


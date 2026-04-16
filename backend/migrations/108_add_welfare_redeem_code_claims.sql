-- Add multi-claim welfare redeem code support.
-- Existing redeem codes remain single-use after backfill.

ALTER TABLE redeem_codes
    ADD COLUMN IF NOT EXISTS max_claims INTEGER NOT NULL DEFAULT 1,
    ADD COLUMN IF NOT EXISTS claimed_count INTEGER NOT NULL DEFAULT 0;

UPDATE redeem_codes
SET max_claims = 1
WHERE max_claims IS NULL OR max_claims <= 0;

UPDATE redeem_codes
SET claimed_count = CASE
    WHEN status = 'used' THEN 1
    ELSE 0
END
WHERE claimed_count IS NULL OR claimed_count = 0;

CREATE INDEX IF NOT EXISTS idx_redeem_codes_type_status ON redeem_codes (type, status);

CREATE TABLE IF NOT EXISTS redeem_code_claims (
    id BIGSERIAL PRIMARY KEY,
    redeem_code_id BIGINT NOT NULL REFERENCES redeem_codes(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount DECIMAL(20,8) NOT NULL,
    claimed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT redeem_code_claims_redeem_code_user_key UNIQUE (redeem_code_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_redeem_code_claims_user_claimed_at ON redeem_code_claims (user_id, claimed_at);
CREATE INDEX IF NOT EXISTS idx_redeem_code_claims_redeem_code_id ON redeem_code_claims (redeem_code_id);

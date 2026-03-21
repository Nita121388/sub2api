-- Add request_log_payloads table for request/response body snapshots.
CREATE TABLE IF NOT EXISTS request_log_payloads (
    id                      BIGSERIAL PRIMARY KEY,
    request_log_id          BIGINT NOT NULL UNIQUE REFERENCES request_logs(id) ON DELETE CASCADE,
    request_body            BYTEA,
    request_body_encoding   VARCHAR(16),
    request_body_bytes      INT,
    request_body_truncated  BOOLEAN NOT NULL DEFAULT FALSE,
    response_body           BYTEA,
    response_body_encoding  VARCHAR(16),
    response_body_bytes     INT,
    response_body_truncated BOOLEAN NOT NULL DEFAULT FALSE,
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_request_log_payloads_request_log_id ON request_log_payloads(request_log_id);
CREATE INDEX IF NOT EXISTS idx_request_log_payloads_created_at ON request_log_payloads(created_at);

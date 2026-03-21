-- Add request_logs table for AI session/log visualization (metadata only).
CREATE TABLE IF NOT EXISTS request_logs (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    api_key_id       BIGINT NOT NULL REFERENCES api_keys(id) ON DELETE CASCADE,
    request_id       VARCHAR(64),
    model            VARCHAR(100) NOT NULL,
    inbound_endpoint VARCHAR(128),
    upstream_endpoint VARCHAR(128),
    method           VARCHAR(8),
    status_code      INT,
    error_code       VARCHAR(64),
    error_message    TEXT,
    input_tokens     INT NOT NULL DEFAULT 0,
    output_tokens    INT NOT NULL DEFAULT 0,
    total_cost       DECIMAL(20, 10) NOT NULL DEFAULT 0,
    stream           BOOLEAN NOT NULL DEFAULT FALSE,
    duration_ms      INT,
    first_token_ms   INT,
    user_agent       VARCHAR(512),
    ip_address       VARCHAR(45),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_request_logs_user_id ON request_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_request_logs_api_key_id ON request_logs(api_key_id);
CREATE INDEX IF NOT EXISTS idx_request_logs_request_id ON request_logs(request_id);
CREATE INDEX IF NOT EXISTS idx_request_logs_model ON request_logs(model);
CREATE INDEX IF NOT EXISTS idx_request_logs_status_code ON request_logs(status_code);
CREATE INDEX IF NOT EXISTS idx_request_logs_created_at ON request_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_request_logs_user_created ON request_logs(user_id, created_at);
CREATE INDEX IF NOT EXISTS idx_request_logs_api_key_created ON request_logs(api_key_id, created_at);

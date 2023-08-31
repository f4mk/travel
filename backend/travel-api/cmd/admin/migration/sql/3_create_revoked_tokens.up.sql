BEGIN;

CREATE TABLE revoked_tokens (
  token_id UUID PRIMARY KEY,
  subject UUID NOT NULL,
  token_version INTEGER NOT NULL,
  issued_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP NOT NULL,
  FOREIGN KEY (subject) REFERENCES users(user_id)
);

COMMIT;
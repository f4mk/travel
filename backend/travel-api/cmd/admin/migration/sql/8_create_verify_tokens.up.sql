BEGIN;

CREATE TABLE verify_tokens (
  token_id TEXT UNIQUE PRIMARY KEY,
  user_id UUID NOT NULL,
  email TEXT NOT NULL,
  issued_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

COMMIT;
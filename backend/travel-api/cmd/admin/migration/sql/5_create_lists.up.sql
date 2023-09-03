BEGIN;

CREATE TABLE lists (
  list_id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  private BOOLEAN NOT NULL DEFAULT TRUE,
  favorite BOOLEAN NOT NULL DEFAULT FALSE,
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  items UUID [] DEFAULT ARRAY [] :: UUID [],
  date_created TIMESTAMP NOT NULL,
  date_updated TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  UNIQUE (user_id, name)
);

COMMIT;
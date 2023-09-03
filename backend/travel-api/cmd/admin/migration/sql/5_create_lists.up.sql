BEGIN;

CREATE TABLE lists (
  list_id UUID PRIMARY KEY,
  user_id UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT,
  private BOOLEAN NOT NULL,
  favorite BOOLEAN NOT NULL,
  completed BOOLEAN NOT NULL,
  items UUID [] DEFAULT ARRAY [] :: UUID [],
  date_created TIMESTAMP NOT NULL,
  date_updated TIMESTAMP NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

COMMIT;
BEGIN;

CREATE TABLE items (
  item_id UUID PRIMARY KEY,
  list_id UUID NOT NULL,
  user_id UUID NOT NULL,
  item_name TEXT NOT NULL,
  description TEXT,
  address TEXT,
  point UUID,
  image_links UUID [] DEFAULT ARRAY [] :: UUID [],
  is_visited BOOLEAN NOT NULL DEFAULT FALSE,
  date_created TIMESTAMP NOT NULL,
  date_updated TIMESTAMP NOT NULL,
  FOREIGN KEY (list_id) REFERENCES lists(list_id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  UNIQUE (list_id, item_name)
);

COMMIT;
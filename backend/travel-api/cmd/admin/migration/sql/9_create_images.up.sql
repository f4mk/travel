BEGIN;

CREATE TYPE image_status AS ENUM ('pending', 'loaded', 'deleted');

CREATE TABLE images (
  image_id UUID PRIMARY KEY,
  list_id UUID NOT NULL,
  user_id UUID NOT NULL,
  item_id UUID NOT NULL,
  private BOOLEAN NOT NULL DEFAULT TRUE,
  description TEXT,
  status image_status NOT NULL,
  date_created TIMESTAMP NOT NULL,
  FOREIGN KEY (list_id) REFERENCES lists(list_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
  UNIQUE (image_id)
);

COMMIT;
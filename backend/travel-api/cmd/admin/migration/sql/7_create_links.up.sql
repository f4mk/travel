BEGIN;

CREATE TABLE links (
  link_id UUID PRIMARY KEY,
  item_id UUID NOT NULL,
  Link_name TEXT,
  url TEXT,
  FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE,
  UNIQUE (item_id, link_name, url)
);

COMMIT;
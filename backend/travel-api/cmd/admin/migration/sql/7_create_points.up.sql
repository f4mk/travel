BEGIN;

CREATE TABLE points (
    point_id UUID PRIMARY KEY,
    item_id UUID NOT NULL,
    location GEOMETRY(POINT, 4326),
    FOREIGN KEY (item_id) REFERENCES items(item_id) ON DELETE CASCADE
);

COMMIT;
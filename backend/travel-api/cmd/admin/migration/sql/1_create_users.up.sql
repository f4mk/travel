BEGIN;

CREATE TABLE users (
	user_id UUID NOT NULL,
	name TEXT NOT NULL,
	email TEXT UNIQUE NOT NULL,
	is_active BOOLEAN NOT NULL DEFAULT false,
	is_deleted BOOLEAN NOT NULL DEFAULT false,
	roles TEXT [] NOT NULL,
	token_version INTEGER NOT NULL,
	password_hash TEXT NOT NULL,
	date_created TIMESTAMP NOT NULL,
	date_updated TIMESTAMP NOT NULL,
	PRIMARY KEY (user_id)
);

COMMIT;
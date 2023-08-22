-- Version: 1.01
-- Description: Create table users
CREATE TABLE users (
	user_id       UUID        NOT NULL,
	name          TEXT        NOT NULL,
	email         TEXT UNIQUE NOT NULL,
	roles         TEXT[]      NOT NULL,
	token_version INTEGER   	NOT NULL,
	password_hash TEXT        NOT NULL,
	date_created  TIMESTAMP   NOT NULL,
	date_updated  TIMESTAMP   NOT NULL,

	PRIMARY KEY (user_id)
);

-- Version: 1.02
-- Description: Create table reset_tokens

CREATE TABLE reset_tokens (
	token_id  		TEXT UNIQUE PRIMARY KEY,
	user_id 			UUID	NOT NULL,
	email 				TEXT NOT NULL, 
	issued_at 		TIMESTAMP NOT NULL,
	expires_at 		TIMESTAMP NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Version: 1.03
-- Description: Create table revoked_tokens
CREATE TABLE revoked_tokens (
    token_id 			UUID PRIMARY KEY,
    subject 			UUID NOT NULL, 
		token_version INTEGER   	NOT NULL,
		issued_at 		TIMESTAMP NOT NULL,
    expires_at 		TIMESTAMP NOT NULL,
    revoked_at 		TIMESTAMP NOT NULL,
    FOREIGN KEY (subject) REFERENCES users(user_id)
);

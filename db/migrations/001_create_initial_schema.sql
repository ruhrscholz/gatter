-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS accounts
(
    account_id BIGSERIAL PRIMARY KEY,
    display_name TEXT,
    fed_username TEXT NOT NULL,
    fed_domain VARCHAR(255) NOT NULL,
        UNIQUE (fed_username, fed_domain),
    id TEXT NOT NULL,
    summary TEXT
);

CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE, -- https://www.rfc-editor.org/rfc/rfc3986
    account_id BIGINT NOT NULL,
        CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES accounts(account_id)
);

CREATE UNLOGGED TABLE IF NOT EXISTS sessions(
    session_id BYTEA PRIMARY KEY,
    user_id INT NOT NULL,
        CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(user_id)
);

CREATE TABLE IF NOT EXISTS applications(
    app_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    scope TEXT NOT NULL,
    website TEXT
);

CREATE TABLE IF NOT EXISTS statuses(
    status_id BIGSERIAL PRIMARY KEY,
    public_id UUID UNIQUE NOT NULL,
    author_id INT,
        CONSTRAINT fk_author FOREIGN KEY(author_id) REFERENCES accounts(account_id),
    created_at TIMESTAMPTZ NOT NULL,
    content TEXT,

    reblogs_count INT,
    favorites_count INT,

    in_reply_to_id INT,
        CONSTRAINT fk_reply_to FOREIGN KEY(in_reply_to_id) REFERENCES statuses(status_id),
    
    reblog INT, -- This will be fun, insert other instances posts into our database
        CONSTRAINT fk_reblog FOREIGN KEY(reblog) REFERENCES statuses(status_id)
);
CREATE UNIQUE INDEX IF NOT EXISTS ind_statuses_public_id on statuses(public_id);

CREATE TABLE IF NOT EXISTS favorites(
    account_id INT PRIMARY KEY,
        CONSTRAINT fk_user FOREIGN KEY(account_id) REFERENCES accounts(account_id),
    status_id BIGINT,
        CONSTRAINT fk_status FOREIGN KEY(status_id) REFERENCES statuses(status_id)
);

-- We don't need a reblog table since we can *hopefully* get that one at request time from status(reblogs_count) and status(reblog)


---- create above / drop below ----

DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS statuses;
DROP INDEX IF EXISTS ind_statuses_public_id;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS applications;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.

CREATE TABLE IF NOT EXISTS accounts
(
    account_id BIGSERIAL PRIMARY KEY,
    display_name TEXT,
    fed_username TEXT NOT NULL;
    fed_domain VARCHAR(255) NOT NULL;
    uri TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(64) NOT NULL,
    password VARCHAR(255) NOT NULL,
    domain VARCHAR(255) UNIQUE, -- https://www.rfc-editor.org/rfc/rfc3986
    account_id BIGINT NOT NULL,
        CONSTRAINT fk_account FOREIGN KEY(account_id) REFERENCES accounts(account_id)
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

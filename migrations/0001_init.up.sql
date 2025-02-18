CREATE TABLE IF NOT EXISTS users (
    username VARCHAR(255) NOT NULL UNIQUE PRIMARY KEY,
    password_hash TEXT NOT NULL,
    coins INT NOT NULL DEFAULT 1000
    );

CREATE TABLE IF NOT EXISTS merch (
    merchname VARCHAR(255) NOT NULL UNIQUE,
    price INT NOT NULL
    );

CREATE TABLE IF NOT EXISTS coin_transactions (
    id BIGSERIAL  PRIMARY KEY,
    from_user_id VARCHAR(255) NOT NULL REFERENCES users (username),
    to_user_id   VARCHAR(255) NOT NULL REFERENCES users (username),
    amount INT NOT NULL
    );

CREATE TABLE IF NOT EXISTS purchases (
    id BIGSERIAL  PRIMARY KEY,
    username VARCHAR(255) NOT NULL REFERENCES users (username),
    merchname VARCHAR(255) NOT NULL REFERENCES merch (merchname),
    quantity INT NOT NULL DEFAULT 1
    );

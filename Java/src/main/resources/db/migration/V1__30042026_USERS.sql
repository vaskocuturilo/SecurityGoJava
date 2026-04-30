CREATE TABLE IF NOT EXISTS users_db
(
    id         BIGINT PRIMARY KEY,
    username   varchar(64) NOT NULL UNIQUE,
    email      varchar(64) NOT NULL UNIQUE,
    password   varchar(64) NOT NULL,
    enabled     boolean,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);
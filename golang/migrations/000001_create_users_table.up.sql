CREATE TABLE IF NOT EXISTS users_db
(
    id         VARCHAR(36)  NOT NULL,
    username   varchar(255) NOT NULL,
    email      varchar(255) NOT NULL,
    password   varchar(255) NOT NULL,
    enabled    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT uq_users_email UNIQUE (email)
);

CREATE INDEX idx_users_email ON users_db (email);
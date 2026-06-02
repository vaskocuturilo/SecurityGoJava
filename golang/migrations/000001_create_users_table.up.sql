CREATE TABLE IF NOT EXISTS users_db
(
    id         VARCHAR(36)  NOT NULL,
    email      VARCHAR(255) NOT NULL,
    password   VARCHAR(255) NOT NULL,
    role       VARCHAR(36)  NOT NULL DEFAULT 'READER',
    enabled    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT chk_users_role CHECK (role IN ('READER', 'CREATE'))
);

CREATE INDEX idx_users_email ON users_db (email);
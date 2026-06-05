CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users_db
(
    id         UUID         NOT NULL DEFAULT gen_random_uuid(),
    username   varchar(255) NOT NULL,
    email      varchar(255) NOT NULL,
    password   varchar(255) NOT NULL,
    role       VARCHAR(36)  NOT NULL DEFAULT 'READER',
    enabled    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,

    CONSTRAINT pk_users PRIMARY KEY (id),
    CONSTRAINT uq_users_email UNIQUE (email),
    CONSTRAINT chk_users_role CHECK (role IN ('READER', 'CREATE'))
);

CREATE INDEX idx_users_email ON users_db (email);

INSERT INTO users_db(id, username, email, password, role, enabled, created_at, updated_at)
VALUES (gen_random_uuid(),
        'admin',
        'admin@test.com',
        '$2a$10$nsG.3gErRUakpvT/63ciiun5eUjhb1c0c3Ud/ZJW58wVAWKzAWPXm',
        'CREATE',
        true,
        current_timestamp,
        current_timestamp)
ON CONFLICT (email) DO NOTHING;
CREATE TABLE IF NOT EXISTS refresh_tokens_db
(
    id         VARCHAR(36)  NOT NULL,
    token      varchar(255) NOT NULL,
    user_id    varchar(36)  NOT NULL,
    expires_at TIMESTAMP    NOT NULL,
    revoked    BOOLEAN      NOT NULL DEFAULT FALSE,

    CONSTRAINT pk_refresh_tokens PRIMARY KEY (id),
    CONSTRAINT uq_refresh_tokens_token UNIQUE (token),
    CONSTRAINT fk_refresh_tokens_user
        FOREIGN KEY (user_id)
            REFERENCES users_db (id)
            ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_token ON refresh_tokens_db (token);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens_db (user_id);
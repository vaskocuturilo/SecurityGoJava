CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_roles
(
    user_id UUID        NOT NULL,
    role    VARCHAR(50) NOT NULL,

    CONSTRAINT pk_user_roles PRIMARY KEY (user_id, role),
    CONSTRAINT fk_user_roles_user
        FOREIGN KEY (user_id)
            REFERENCES users_db (id)
            ON DELETE CASCADE,
    CONSTRAINT chk_user_roles_role
        CHECK (role IN ('READER', 'CREATE'))
);

CREATE INDEX idx_user_roles_user_id ON user_roles (user_id);

INSERT INTO user_roles (user_id, role)
SELECT id, 'CREATE'
FROM users_db
WHERE email = 'admin@test.com'
ON CONFLICT DO NOTHING;

INSERT INTO user_roles (user_id, role)
SELECT id, 'READER'
FROM users_db
WHERE email = 'admin@test.com'
ON CONFLICT DO NOTHING;
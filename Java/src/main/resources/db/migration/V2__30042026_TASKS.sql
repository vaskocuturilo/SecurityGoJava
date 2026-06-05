CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS tasks_db
(
    id          UUID NOT NULL DEFAULT gen_random_uuid(),
    title       VARCHAR(255),
    description VARCHAR(255),

    CONSTRAINT pk_tasks PRIMARY KEY (id)
);

CREATE INDEX idx_tasks_title ON tasks_db (title);
CREATE INDEX idx_tasks_description ON tasks_db (description);
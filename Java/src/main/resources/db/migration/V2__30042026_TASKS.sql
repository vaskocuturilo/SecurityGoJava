CREATE TABLE IF NOT EXISTS tasks_db
(
    id          VARCHAR(36) NOT NULL,
    title       VARCHAR(255),
    description VARCHAR(255),

    CONSTRAINT pk_tasks PRIMARY KEY (id)
);

CREATE INDEX idx_tasks_title ON tasks_db (title);
CREATE INDEX idx_tasks_description ON tasks_db (description);
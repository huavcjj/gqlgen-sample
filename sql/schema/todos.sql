CREATE TABLE IF NOT EXISTS users
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS todos
(
    id      TEXT PRIMARY KEY,
    text    TEXT    NOT NULL,
    done    BOOLEAN NOT NULL,
    user_id TEXT    NOT NULL,
    FOREIGN KEY (user_id) REFERENCES USERS (id)
)

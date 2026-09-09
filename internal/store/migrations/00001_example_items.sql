-- +goose Up
CREATE TABLE example_items (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    done       INTEGER NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX example_items_done_idx ON example_items (done);

-- +goose Down
DROP INDEX example_items_done_idx;

DROP TABLE example_items;

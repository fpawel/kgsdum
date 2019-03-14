PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS work
(
  work_id    INTEGER   NOT NULL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL UNIQUE DEFAULT (DATETIME('now')),
  name       TEXT      NOT NULL
);

CREATE TABLE IF NOT EXISTS entry
(
  entry_id   INTEGER   NOT NULL PRIMARY KEY,
  work_id    INTEGER   NOT NULL,
  created_at TIMESTAMP NOT NULL UNIQUE DEFAULT (DATETIME('now')),
  level      INTEGER   NOT NULL,
  message    TEXT      NOT NULL,
  FOREIGN KEY (work_id) REFERENCES work (work_id) ON DELETE CASCADE
);

CREATE VIEW IF NOT EXISTS last_work AS
SELECT * FROM work
ORDER BY created_at DESC
LIMIT 1;

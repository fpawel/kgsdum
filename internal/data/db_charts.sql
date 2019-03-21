PRAGMA foreign_keys = ON;
PRAGMA encoding = 'UTF-8';

CREATE TABLE IF NOT EXISTS bucket
(
  bucket_id  INTEGER NOT NULL PRIMARY KEY,
  created_at REAL    NOT NULL UNIQUE DEFAULT (julianday('now')),
  updated_at REAL    NOT NULL        DEFAULT (julianday('now')),
  name       TEXT    NOT NULL
);

CREATE TABLE IF NOT EXISTS series
(
  bucket_id INTEGER NOT NULL,
  addr      INTEGER NOT NULL CHECK (addr > 0),
  var       INTEGER NOT NULL CHECK (var >= 0),
  stored_at REAL    NOT NULL,
  value     REAL    NOT NULL,
  FOREIGN KEY (bucket_id) REFERENCES bucket (bucket_id)
    ON DELETE CASCADE
);

CREATE TRIGGER IF NOT EXISTS trigger_bucket_updated_at
  AFTER INSERT
  ON series
  FOR EACH ROW
  BEGIN
    UPDATE bucket
    SET updated_at = julianday('now')
    WHERE bucket.bucket_id = new.bucket_id;
  END;

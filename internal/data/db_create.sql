CREATE TABLE IF NOT EXISTS party
(
  party_id     INTEGER PRIMARY KEY      NOT NULL,
  created_at   TIMESTAMP WITH TIME ZONE NOT NULL                     DEFAULT (datetime('now')) UNIQUE,
  product_type TEXT                     NOT NULL                     DEFAULT '00.01',
  pgs1         REAL                     NOT NULL CHECK ( pgs1 >= 0 ) DEFAULT 0,
  pgs2         REAL                     NOT NULL CHECK ( pgs2 >= 0 ) DEFAULT 4,
  pgs3         REAL                     NOT NULL CHECK ( pgs3 >= 0 ) DEFAULT 7.5,
  pgs4         REAL                     NOT NULL CHECK ( pgs4 >= 0 ) DEFAULT 12
);

CREATE TABLE IF NOT EXISTS product
(
  product_id    INTEGER PRIMARY KEY      NOT NULL,
  party_id      INTEGER                  NOT NULL,
  created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (datetime('now')) UNIQUE,
  serial_number TEXT                     NOT NULL CHECK (serial_number <> '' ),
  addr          SMALLINT                 NOT NULL CHECK (addr > 0),
  production    BOOLEAN                  NOT NULL DEFAULT FALSE,

  c_norm1       REAL,
  c_norm2       REAL,
  c_norm3       REAL,
  c_norm4       REAL,

  c_plus1       REAL,
  c_plus2       REAL,
  c_plus3       REAL,
  c_plus4       REAL,

  c_minus1      REAL,
  c_minus2      REAL,
  c_minus3      REAL,
  c_minus4      REAL,

  UNIQUE (party_id, addr),
  UNIQUE (party_id, serial_number),
  FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);


CREATE VIEW IF NOT EXISTS last_party AS
SELECT *
FROM party
ORDER BY created_at DESC
LIMIT 1;

CREATE VIEW IF NOT EXISTS last_party_id AS
SELECT party_id
FROM party
ORDER BY created_at DESC
LIMIT 1;

CREATE VIEW IF NOT EXISTS last_party_product AS
SELECT *
FROM product
WHERE party_id IN (SELECT * FROM last_party_id)
ORDER BY created_at;

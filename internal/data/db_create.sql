CREATE TABLE IF NOT EXISTS party
(
  party_id          SERIAL PRIMARY KEY       NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL                        DEFAULT CURRENT_TIMESTAMP,
  product_type      TEXT                     NOT NULL                        DEFAULT '00.01',
  pgs_beg           REAL                     NOT NULL CHECK ( pgs_beg >= 0 ) DEFAULT 0,
  pgs_mid           REAL                     NOT NULL CHECK ( pgs_mid >= 0 ) DEFAULT 50,
  pgs_end           REAL                     NOT NULL CHECK ( pgs_end >= 0 ) DEFAULT 100,
  temperature_norm  REAL                     NOT NULL                        DEFAULT 20,
  temperature_plus  REAL                     NOT NULL                        DEFAULT 60,
  temperature_minus REAL                     NOT NULL                        DEFAULT -30
);

CREATE TABLE IF NOT EXISTS product
(
  product_id           SERIAL PRIMARY KEY NOT NULL,
  party_id             INTEGER            NOT NULL,
  serial_number        TEXT               NOT NULL CHECK (serial_number <> '' ),
  addr                 SMALLINT           NOT NULL CHECK (addr > 0),
  production           BOOLEAN            NOT NULL DEFAULT FALSE,
  connection_error     TEXT,
  connection_beg_norm  REAL,
  connection_mid_norm  REAL,
  connection_end_norm  REAL,
  connection_beg_minus REAL,
  connection_mid_minus REAL,
  connection_end_minus REAL,
  connection_beg_plus  REAL,
  connection_mid_plus  REAL,
  connection_end_plus  REAL,
  UNIQUE (party_id, addr),
  UNIQUE (party_id, serial_number),
  FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);

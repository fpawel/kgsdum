CREATE TABLE IF NOT EXISTS party
(
  party_id          SERIAL PRIMARY KEY       NOT NULL,
  created_at        TIMESTAMP WITH TIME ZONE NOT NULL                        DEFAULT CURRENT_TIMESTAMP UNIQUE,
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
  product_id              SERIAL PRIMARY KEY       NOT NULL,
  party_id                INTEGER                  NOT NULL,
  created_at              TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP UNIQUE,
  serial_number           TEXT                     NOT NULL CHECK (serial_number <> '' ),
  addr                    SMALLINT                 NOT NULL CHECK (addr > 0),
  production              BOOLEAN                  NOT NULL DEFAULT FALSE,
  connection_error        TEXT,
  concentration_beg_norm  REAL,
  concentration_mid_norm  REAL,
  concentration_end_norm  REAL,
  concentration_beg_minus REAL,
  concentration_mid_minus REAL,
  concentration_end_minus REAL,
  concentration_beg_plus  REAL,
  concentration_mid_plus  REAL,
  concentration_end_plus  REAL,
  UNIQUE (party_id, addr),
  UNIQUE (party_id, serial_number),
  FOREIGN KEY (party_id) REFERENCES party (party_id) ON DELETE CASCADE
);

CREATE OR REPLACE VIEW last_party AS
SELECT *
FROM party
ORDER BY created_at DESC
LIMIT 1;

CREATE OR REPLACE VIEW last_party_id AS
SELECT party_id
FROM party
ORDER BY created_at DESC
LIMIT 1;

CREATE OR REPLACE VIEW last_party_product AS
SELECT *
FROM product
WHERE party_id IN (SELECT * FROM last_party_id)
ORDER BY created_at;


CREATE OR REPLACE FUNCTION err_percent1(float) RETURNS float AS $$
SELECT ROUND(($1 * 100)::numeric,1)::float;
$$ language SQL IMMUTABLE;

CREATE OR REPLACE VIEW product_info AS
  WITH q1 AS (
    SELECT product.*,

           concentration_beg_norm - pgs_beg   AS d_beg_norm,
           concentration_mid_norm - pgs_mid  AS d_mid_norm,
           concentration_end_norm - pgs_end  AS d_end_norm,

           concentration_beg_minus - pgs_beg  AS d_beg_minus,
           concentration_mid_minus - pgs_mid  AS d_mid_minus,
           concentration_end_minus - pgs_end  AS d_end_minus,

           concentration_beg_plus - pgs_beg   AS d_beg_plus,
           concentration_mid_plus - pgs_mid   AS d_mid_plus,
           concentration_end_plus - pgs_end   AS d_end_plus,

           (0.1 + 0.12 * party.pgs_beg)      AS err_beg_limit,
           (0.1 + 0.12 * party.pgs_mid)      AS err_mid_limit,
           (0.1 + 0.12 * party.pgs_end)      AS err_end_limit
    FROM product
           INNER JOIN party on product.party_id = party.party_id
    )
    SELECT q1.*,

           err_percent1( d_beg_norm / err_beg_limit)  AS err_beg_norm_percent,
           err_percent1( d_mid_norm / err_mid_limit)  AS err_mid_norm_percent,
           err_percent1( d_end_norm / err_end_limit)  AS err_end_norm_percent,

           err_percent1( d_beg_minus / err_beg_limit)  AS err_beg_minus_percent,
           err_percent1(d_mid_minus / err_mid_limit)  AS err_mid_minus_percent,
           err_percent1( d_end_minus / err_end_limit)  AS err_end_minus_percent,

           err_percent1( d_beg_plus / err_beg_limit)  AS err_beg_plus_percent,
           err_percent1( d_mid_plus / err_mid_limit)  AS err_mid_plus_percent,
           err_percent1( d_end_plus / err_end_limit)  AS err_end_plus_percent,

           abs(d_beg_norm) < err_beg_limit  AS ok_beg_norm,
           abs(d_mid_norm) < err_mid_limit  AS ok_mid_norm,
           abs(d_end_norm) < err_end_limit  AS ok_end_norm,

           abs(d_beg_minus) < err_beg_limit AS ok_beg_minus,
           abs(d_mid_minus) < err_mid_limit AS ok_mid_minus,
           abs(d_end_minus) < err_end_limit AS ok_end_minus,

           abs(d_beg_plus) < err_beg_limit  AS ok_beg_plus,
           abs(d_mid_plus) < err_mid_limit  AS ok_mid_plus,
           abs(d_end_plus) < err_end_limit  AS ok_end_plus
    FROM q1;

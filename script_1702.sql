CREATE TABLE IF NOT EXISTS company
(
  id          SERIAL NOT NULL
    CONSTRAINT group_pkey
    PRIMARY KEY,
  name        VARCHAR(200),
  description TEXT,
  inn         VARCHAR(20),
  kpp         VARCHAR(20),
  okpo        VARCHAR(20),
  ogrn        VARCHAR(20),
  email_index VARCHAR(10),
  street_id   INTEGER,
  contact     VARCHAR(20),
  confirm     BOOLEAN DEFAULT FALSE,
  pub         BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS group_name_uindex
  ON company (name);

CREATE TABLE IF NOT EXISTS client
(
  id         SERIAL                                          NOT NULL
    CONSTRAINT user_pkey
    PRIMARY KEY,
  email      VARCHAR(200)                                    NOT NULL,
  pass       VARCHAR(100)                                    NOT NULL,
  confirmed  BOOLEAN DEFAULT FALSE,
  first_name VARCHAR(100) DEFAULT 'tmp' :: CHARACTER VARYING NOT NULL,
  last_name  VARCHAR(100) DEFAULT 'tmp' :: CHARACTER VARYING NOT NULL,
  pub        BOOLEAN DEFAULT FALSE,
  verified   BOOLEAN DEFAULT FALSE,
  avatar     TEXT,
  country_id INTEGER
);

CREATE UNIQUE INDEX IF NOT EXISTS user_email_uindex
  ON client (email);

CREATE TABLE IF NOT EXISTS document
(
  id           SERIAL                NOT NULL
    CONSTRAINT document_pkey
    PRIMARY KEY,
  name         VARCHAR(500),
  client_id    INTEGER
    CONSTRAINT document_user_id_fk
    REFERENCES client
    ON DELETE SET NULL,
  public       BOOLEAN DEFAULT FALSE NOT NULL,
  template     BOOLEAN DEFAULT FALSE NOT NULL,
  sent         BOOLEAN DEFAULT FALSE,
  last_updated TIMESTAMP,
  created      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS client_company
(
  client_id    INTEGER NOT NULL
    CONSTRAINT client_company_client_id_fk
    REFERENCES client
    ON DELETE CASCADE,
  company_id   INTEGER NOT NULL
    CONSTRAINT client_company_company_id_fk
    REFERENCES company,
  gr_admin     BOOLEAN DEFAULT FALSE,
  gr_invite    BOOLEAN DEFAULT FALSE,
  gr_kick      BOOLEAN DEFAULT FALSE,
  gr_read      BOOLEAN DEFAULT FALSE,
  gr_write     BOOLEAN DEFAULT FALSE,
  gr_update    BOOLEAN DEFAULT FALSE,
  gr_delete    BOOLEAN DEFAULT FALSE,
  responsible  BOOLEAN DEFAULT FALSE,
  confirm      BOOLEAN,
  company_conf BOOLEAN,
  CONSTRAINT client_company_client_id_company_id_pk
  PRIMARY KEY (client_id, company_id)
);

CREATE TABLE IF NOT EXISTS block
(
  id        SERIAL NOT NULL
    CONSTRAINT block_pkey
    PRIMARY KEY,
  name      VARCHAR(1000),
  content   TEXT,
  shared    BOOLEAN DEFAULT FALSE,
  date      TIMESTAMP,
  client_id INTEGER DEFAULT 1
    CONSTRAINT block_client_id_fk
    REFERENCES client
    ON UPDATE CASCADE ON DELETE CASCADE,
  fts       TSVECTOR,
  copy      INTEGER
);

CREATE OR REPLACE FUNCTION create_ts_block_data()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.fts := to_tsvector('ru', NEW.name) || to_tsvector('ru', NEW.content);
  RETURN NEW;
END;
$$;

CREATE TRIGGER create_ts_data
  BEFORE INSERT OR UPDATE
  ON block
  FOR EACH ROW
EXECUTE PROCEDURE create_ts_block_data();

CREATE OR REPLACE FUNCTION add_old_block()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO block_history (name, content, date, block_id, doc, parent, ord, deleted)
    SELECT
      old.name,
      old.content,
      current_timestamp,
      old.id,
      doc_id,
      parent_block,
      ord,
      FALSE
    FROM doc_block
      JOIN block b ON doc_block.block_id = b.id
      JOIN document d2 ON doc_block.doc_id = d2.id
    WHERE block_id = old.id AND b.client_id = d2.client_id;
  RETURN new;
END;
$$;

CREATE TRIGGER trigger_update_block
  BEFORE UPDATE
  ON block
  FOR EACH ROW
EXECUTE PROCEDURE add_old_block();

CREATE TABLE IF NOT EXISTS doc_block
(
  doc_id       INTEGER NOT NULL
    CONSTRAINT doc_block_document_id_fk
    REFERENCES document
    ON DELETE CASCADE,
  block_id     INTEGER NOT NULL
    CONSTRAINT doc_block_block_id_fk
    REFERENCES block
    ON DELETE CASCADE,
  ord          INTEGER,
  parent_block INTEGER
    CONSTRAINT doc_block_parent__fk
    REFERENCES block
    ON DELETE CASCADE,
  CONSTRAINT doc_block_block_id_pk
  PRIMARY KEY (doc_id, block_id)
);

CREATE OR REPLACE FUNCTION refactor_ord_insert()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE doc_block
  SET ord = ord + 1
  WHERE doc_block.block_id != new.block_id AND ord >= new.ord AND doc_id = new.doc_id AND
        parent_block = new.parent_block;
  RETURN new;
END;
$$;

CREATE TRIGGER trigger_insert
  AFTER INSERT
  ON doc_block
  FOR EACH ROW
EXECUTE PROCEDURE refactor_ord_insert();

CREATE OR REPLACE FUNCTION refactor_ord_delete()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  UPDATE doc_block
  SET ord = ord - 1
  WHERE ord > old.ord AND doc_id = old.doc_id AND
        parent_block = old.parent_block;
  RETURN new;
END;
$$;

CREATE TRIGGER trigger_delete
  AFTER DELETE
  ON doc_block
  FOR EACH ROW
EXECUTE PROCEDURE refactor_ord_delete();

CREATE TABLE IF NOT EXISTS company_doc
(
  company_id INTEGER NOT NULL
    CONSTRAINT group_block_group_id_fk
    REFERENCES company
    ON DELETE CASCADE,
  doc_id     INTEGER NOT NULL
    CONSTRAINT company_doc_document_id_fk
    REFERENCES document,
  CONSTRAINT group_block_block_id_pk
  PRIMARY KEY (doc_id, company_id)
);

CREATE TABLE IF NOT EXISTS block_history
(
  id       SERIAL                NOT NULL
    CONSTRAINT block_history_pkey
    PRIMARY KEY,
  name     VARCHAR(1000),
  content  TEXT,
  date     TIMESTAMP,
  block_id INTEGER
    CONSTRAINT block_history_block_id_fk
    REFERENCES block
    ON UPDATE CASCADE ON DELETE CASCADE,
  ord      INTEGER,
  parent   INTEGER,
  doc      INTEGER,
  deleted  BOOLEAN DEFAULT FALSE NOT NULL
);

CREATE TABLE IF NOT EXISTS client_confirm
(
  client_id    INTEGER NOT NULL
    CONSTRAINT client_confirm_client_id_pk
    PRIMARY KEY
    CONSTRAINT user_confirm_user_id_fk
    REFERENCES client
    ON DELETE CASCADE,
  uid          VARCHAR(100),
  date_exp     TIMESTAMP,
  date_confirm TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS user_confirm_uid_uindex
  ON client_confirm (uid);

CREATE TABLE IF NOT EXISTS request_company
(
  uid       VARCHAR(80) NOT NULL
    CONSTRAINT request_company_uid_pk
    PRIMARY KEY,
  client_id INTEGER     NOT NULL
    CONSTRAINT request_company_client_id_fk
    REFERENCES client
    ON DELETE CASCADE,
  name      VARCHAR(50) NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS request_company_uid_uindex
  ON request_company (uid);

CREATE UNIQUE INDEX IF NOT EXISTS request_company_name_uindex
  ON request_company (name);

CREATE TABLE IF NOT EXISTS country
(
  id   SERIAL NOT NULL
    CONSTRAINT country_pkey
    PRIMARY KEY,
  name VARCHAR(200)
);

CREATE UNIQUE INDEX IF NOT EXISTS country_name_uindex
  ON country (name);

ALTER TABLE client
  ADD CONSTRAINT client_country_id_fk
FOREIGN KEY (country_id) REFERENCES country;

CREATE TABLE IF NOT EXISTS city
(
  id         SERIAL       NOT NULL
    CONSTRAINT city_id_pk
    PRIMARY KEY,
  name       VARCHAR(200) NOT NULL,
  country_id INTEGER
    CONSTRAINT city_country_id_fk
    REFERENCES country
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS city_name_uindex
  ON city (name);

CREATE TABLE IF NOT EXISTS street
(
  id      SERIAL       NOT NULL
    CONSTRAINT street_pkey
    PRIMARY KEY,
  name    VARCHAR(200) NOT NULL,
  city_id INTEGER
    CONSTRAINT street_city_id_fk
    REFERENCES city
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS street_name_uindex
  ON street (name);

ALTER TABLE company
  ADD CONSTRAINT company_street_id_fk
FOREIGN KEY (street_id) REFERENCES street
ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS recieve_document
(
  client_id  INTEGER
    CONSTRAINT recieve_document_client_id_fk
    REFERENCES client
    ON UPDATE CASCADE ON DELETE CASCADE,
  document   INTEGER
    CONSTRAINT recieve_document_document_id_fk
    REFERENCES document
    ON UPDATE CASCADE ON DELETE CASCADE,
  company_id INTEGER
    CONSTRAINT recieve_document_company_id_fk
    REFERENCES company
    ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE OR REPLACE FUNCTION get_full_block(block_v INTEGER)
  RETURNS TABLE(id INTEGER, parent_id INTEGER, level INTEGER, content TEXT, name CHARACTER VARYING, ord INTEGER)
LANGUAGE SQL
AS $$
WITH RECURSIVE get_all_blocks(id, parent_id, level, content, name, ord) AS (
  SELECT
    block.id,
    db.parent_block,
    0,
    block.content,
    block.name,
    db.ord
  FROM block
    JOIN doc_block db ON block.id = db.block_id
    JOIN document d2 ON db.doc_id = d2.id
  WHERE d2.client_id = block.client_id AND db.block_id = block_v
  UNION ALL
  SELECT
    block.id,
    db.parent_block,
    level + 1,
    block.content,
    block.name,
    db.ord
  FROM block
    JOIN doc_block db ON block.id = db.block_id
    JOIN get_all_blocks rt ON rt.id = db.parent_block
)
SELECT *
FROM get_all_blocks
ORDER BY id;
$$;

CREATE OR REPLACE FUNCTION get_all_blocks(doc INTEGER)
  RETURNS TABLE(id INTEGER, parent_id INTEGER, level INTEGER, content TEXT, name CHARACTER VARYING, ord INTEGER)
LANGUAGE SQL
AS $$
WITH RECURSIVE get_all_blocks(id, parent_id, level, content, name, ord) AS
(SELECT
   block.id,
   db.parent_block,
   0,
   block.content,
   block.name,
   db.ord
 FROM block
   JOIN doc_block db ON block.id = db.block_id
 WHERE db.doc_id = doc AND db.parent_block IS NULL
 UNION ALL
 SELECT
   block.id,
   db.parent_block,
   level + 1,
   block.content,
   block.name,
   db.ord
 FROM block
   JOIN doc_block db ON block.id = db.block_id
   JOIN get_all_blocks rt ON rt.id = db.parent_block
)
SELECT *
FROM get_all_blocks;
$$;

CREATE OR REPLACE FUNCTION delete_block(doc INTEGER, block INTEGER)
  RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  DELETE FROM doc_block
  WHERE doc_id = doc AND block_id IN (SELECT id
                                      FROM get_full_block(block));
END;
$$;

CREATE OR REPLACE FUNCTION full_delete_block(del_block INTEGER)
  RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO block_history (name, content, date, ord, parent, doc, deleted)
    SELECT
      block.name,
      block.content,
      block.date,
      db.ord,
      db.parent_block,
      db.doc_id,
      TRUE
    FROM block
      JOIN doc_block db ON block.id = db.block_id
      JOIN document d2 ON db.doc_id = d2.id
    WHERE block.client_id = d2.client_id AND block.id IN (SELECT id
                                                          FROM get_full_block(del_block));

  WITH tmp AS (
    DELETE FROM doc_block
    WHERE block_id IN (SELECT id
                       FROM get_full_block(del_block))
    RETURNING block_id
  )

  DELETE FROM block
  WHERE id IN (SELECT block_id
               FROM TMP);
END;
$$;

CREATE OR REPLACE FUNCTION copy_doc(doc INTEGER, usr INTEGER)
  RETURNS INTEGER
LANGUAGE plpgsql
AS $$
DECLARE
  new_doc INT;
BEGIN
  INSERT INTO document (name, client_id) SELECT
                                           name,
                                           usr
                                         FROM document
                                         WHERE id = doc
  RETURNING id
    INTO new_doc;
  INSERT INTO block (name, content, date, client_id, copy) SELECT
                                                             b.name,
                                                             b.content,
                                                             current_timestamp,
                                                             usr,
                                                             b.id
                                                           FROM block b
                                                             JOIN doc_block db ON b.id = db.block_id
                                                           WHERE db.doc_id = doc;
  INSERT INTO doc_block (doc_id, block_id, ord) (SELECT
                                                   new_doc,
                                                   b2.id,
                                                   db.ord
                                                 FROM doc_block db
                                                   JOIN block b2 ON db.block_id = b2.copy
                                                 WHERE db.doc_id = doc AND db.parent_block ISNULL
  );
  INSERT INTO doc_block (doc_id, block_id, ord, parent_block) (SELECT
                                                                 new_doc,
                                                                 b2.id,
                                                                 db.ord,
                                                                 b3.id
                                                               FROM doc_block db
                                                                 JOIN block b2 ON db.block_id = b2.copy
                                                                 JOIN block b3 ON db.parent_block = b3.copy
                                                               WHERE db.doc_id = doc
  );
  UPDATE block
  SET copy = NULL
  WHERE copy IS NOT NULL;
  RETURN new_doc;
END;
$$;

CREATE OR REPLACE FUNCTION delete_doc(doc INTEGER)
  RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  WITH tmp AS (
    DELETE FROM doc_block
    WHERE doc_id = doc
    RETURNING block_id)
  DELETE FROM block
  WHERE id IN (SELECT block_id
               FROM tmp);
  DELETE FROM document
  WHERE id = doc;
END;
$$;

CREATE OR REPLACE FUNCTION history_delete_block()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  INSERT INTO block_history (name, content, date, doc, parent, ord, deleted) SELECT
                                                                               old.name,
                                                                               old.content,
                                                                               current_timestamp,
                                                                               doc_id,
                                                                               parent_block,
                                                                               ord,
                                                                               TRUE
                                                                             FROM doc_block
                                                                               JOIN block b ON doc_block.block_id = b.id
                                                                               JOIN document d2
                                                                                 ON doc_block.doc_id = d2.id
                                                                             WHERE block_id = old.id AND
                                                                                   b.client_id = d2.client_id;
  DELETE FROM block_history
  WHERE block_id = old.id;
  RETURN new;
END;
$$;

CREATE OR REPLACE FUNCTION delete_block_refer(doc INTEGER, block INTEGER)
  RETURNS VOID
LANGUAGE plpgsql
AS $$
BEGIN
  DELETE FROM doc_block
  WHERE doc_id = doc AND block_id IN (SELECT id
                                      FROM get_full_block(block));
END;
$$;

-- missing source code for set_limit
;

-- missing source code for show_limit
;

-- missing source code for show_trgm
;

-- missing source code for similarity
;

-- missing source code for similarity_op
;

-- missing source code for similarity_dist
;

-- missing source code for gtrgm_in
;

-- missing source code for gtrgm_out
;

-- missing source code for gtrgm_consistent
;

-- missing source code for gtrgm_distance
;

-- missing source code for gtrgm_compress
;

-- missing source code for gtrgm_decompress
;

-- missing source code for gtrgm_penalty
;

-- missing source code for gtrgm_picksplit
;

-- missing source code for gtrgm_union
;

-- missing source code for gtrgm_same
;

-- missing source code for gin_extract_value_trgm
;

-- missing source code for gin_extract_query_trgm
;

-- missing source code for gin_trgm_consistent
;

-- missing source code for dxsyn_init
;

-- missing source code for dxsyn_lexize
;

CREATE OR REPLACE FUNCTION create_ts_data()
  RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.fts := to_tsvector('ru', NEW.name) || to_tsvector('ru', NEW.content);
  RETURN NEW;
END;
$$;

CREATE OPERATOR % ( PROCEDURE = similarity_op, LEFTARG = TEXT, RIGHTARG = TEXT );

CREATE OPERATOR <-> ( PROCEDURE = similarity_dist, LEFTARG = TEXT, RIGHTARG = TEXT );



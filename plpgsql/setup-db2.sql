\ir ./setup-db-primary.sql;

CREATE OR REPLACE VIEW "ViewSum1" AS (
  SELECT
    "Table1"."ID" || ':' || "Table2"."ID" AS "ID",
    "Table1"."ID" AS "Table1ID",
    "Table2"."ID" AS "Table2ID",
    "Table1"."I" AS "Table1I",
    "Table2"."I" AS "Table2I",
    "Table1"."Num1" AS "Table1Num1",
    "Table2"."Num3" AS "Table2Num3",
    "Table1"."Num1" + "Table2"."Num3" AS "Sum13"
  FROM "Table1", "Table2"
);

CREATE OR REPLACE FUNCTION process_viewsum1()
  RETURNS trigger
  LANGUAGE 'plpgsql'
AS $$
DECLARE
  r RECORD;
BEGIN
IF TG_OP = 'UPDATE' THEN
  IF (NEW."Table1I" <> OLD."Table1I") THEN
    FOR r IN (
      SELECT
        'Table1' AS source,
        lower(TG_OP) AS method,
        row_to_json(t) AS result,
        TRUE AS primary
      FROM (
        SELECT
          NEW."Table1ID" AS "ID",
          NEW."Table1I" AS "I"
      ) t
    ) LOOP
      PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
    END LOOP;
  END IF;
  IF (NEW."Table2I" <> OLD."Table2I") THEN
    FOR r IN (
      SELECT
        'Table1' AS source,
        lower(TG_OP) AS method,
        row_to_json(t) AS result,
        TRUE AS primary
      FROM (
        SELECT
          NEW."Table2ID" AS "ID",
          NEW."Table2I" AS "I"
      ) t
    ) LOOP
      PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
    END LOOP;
  END IF;
  IF (NEW."Table1Num1" <> OLD."Table1Num1") THEN
    FOR r IN (
      SELECT
        'Table1' AS source,
        lower(TG_OP) AS method,
        row_to_json(t) AS result,
        TRUE AS primary
      FROM (
        SELECT
          NEW."Table1ID" AS "ID",
          NEW."Table1I" AS "I",
          NEW."Table1Num1" AS "Num1"
      ) t
    ) LOOP
      PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
    END LOOP;
  END IF;
  IF (NEW."Table2Num3" <> OLD."Table2Num3") THEN
    FOR r IN (
      SELECT
        'Table2' AS source,
        lower(TG_OP) AS method,
        row_to_json(t) AS result,
        TRUE AS primary
      FROM (
        SELECT
          NEW."Table2ID" AS "ID",
          NEW."Table2I" AS "I",
          NEW."Table2Num3" AS "Num3"
      ) t
    ) LOOP
      PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
    END LOOP;
  END IF;
  RETURN NEW;
END IF;
RETURN NULL;
END;
$$;

DROP TRIGGER IF EXISTS "ViewSum1_process" ON "ViewSum1";
CREATE TRIGGER "ViewSum1_process"
INSTEAD OF INSERT OR UPDATE OR DELETE ON "ViewSum1"
FOR EACH ROW
EXECUTE PROCEDURE process_viewsum1();

CREATE OR REPLACE FUNCTION notify_from_table1_viewsum1_before()
  RETURNS TRIGGER
  LANGUAGE plpgsql
AS $$
DECLARE
  r RECORD;
BEGIN
IF (TG_OP = 'DELETE') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table1ID"=OLD."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN OLD;
ELSIF (TG_OP = 'UPDATE') THEN
  RETURN NEW;
ELSIF (TG_OP = 'INSERT') THEN
  RETURN NEW;
END IF;
RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION notify_from_table1_viewsum1_after()
  RETURNS TRIGGER
  LANGUAGE plpgsql
AS $$
DECLARE
  r RECORD;
BEGIN
IF (TG_OP = 'DELETE') THEN
  RETURN OLD;
ELSIF (TG_OP = 'UPDATE') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table1ID"=NEW."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN NEW;
ELSIF (TG_OP = 'INSERT') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table1ID"=NEW."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN NEW;
END IF;
RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION notify_from_table2_viewsum1_before()
  RETURNS TRIGGER
  LANGUAGE plpgsql
AS $$
DECLARE
  r RECORD;
BEGIN
IF (TG_OP = 'DELETE') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table2ID"=OLD."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN OLD;
ELSIF (TG_OP = 'UPDATE') THEN
  RETURN NEW;
ELSIF (TG_OP = 'INSERT') THEN
  RETURN NEW;
END IF;
RETURN NULL;
END;
$$;

CREATE OR REPLACE FUNCTION notify_from_table2_viewsum1_after()
  RETURNS TRIGGER
  LANGUAGE plpgsql
AS $$
DECLARE
  r RECORD;
BEGIN
IF (TG_OP = 'DELETE') THEN
  RETURN OLD;
ELSIF (TG_OP = 'UPDATE') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table2ID"=NEW."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN NEW;
ELSIF (TG_OP = 'INSERT') THEN
  FOR r IN (
    SELECT
      'ViewSum1' AS source,
      lower(TG_OP) AS method,
      row_to_json(t) AS result
    FROM (
      SELECT *
        FROM "ViewSum1"
        WHERE "Table2ID"=NEW."ID"
    ) t
  ) LOOP
    PERFORM pg_notify('jsonrpc', row_to_json(r)::text);
  END LOOP;
  RETURN NEW;
END IF;
RETURN NULL;
END;
$$;

DROP TRIGGER IF EXISTS "Table1_notify_viewsum1_before" ON "Table1";
CREATE TRIGGER "Table1_notify_viewsum1_before"
  BEFORE INSERT OR UPDATE OR DELETE
  ON "Table1"
  FOR EACH ROW
  EXECUTE PROCEDURE notify_from_table1_viewsum1_before();

DROP TRIGGER IF EXISTS "Table1_notify_viewsum1_after" ON "Table1";
CREATE TRIGGER "Table1_notify_viewsum1_after"
  AFTER INSERT OR UPDATE OR DELETE
  ON "Table1"
  FOR EACH ROW
  EXECUTE PROCEDURE notify_from_table1_viewsum1_after();

DROP TRIGGER IF EXISTS "Table2_notify_viewsum1_before" ON "Table2";
CREATE TRIGGER "Table2_notify_viewsum1_before"
  BEFORE INSERT OR UPDATE OR DELETE
  ON "Table2"
  FOR EACH ROW
  EXECUTE PROCEDURE notify_from_table2_viewsum1_before();

DROP TRIGGER IF EXISTS "Table2_notify_viewsum1_after" ON "Table2";
CREATE TRIGGER "Table2_notify_viewsum1_after"
  AFTER INSERT OR UPDATE OR DELETE
  ON "Table2"
  FOR EACH ROW
  EXECUTE PROCEDURE notify_from_table2_viewsum1_after();

/*
  Delete the databases before running this function.
  The IDs MUST be equal everywhere
*/
CREATE OR REPLACE FUNCTION goahead(i BIGINT)
  RETURNS VOID
  LANGUAGE 'plpgsql'
AS $$
DECLARE
  c111 double precision=CEIL(RANDOM() * 1000) / 10;
  c123 double precision=CEIL(RANDOM() * 1000) / 10;
  c211 double precision=CEIL(RANDOM() * 1000) / 10;
  c223 double precision=CEIL(RANDOM() * 1000) / 10;
BEGIN
RAISE NOTICE 'I=% c111=% c123=%', i, c111, c123;
RAISE NOTICE 'I=% c211=% c223=%', i, c211, c223;
UPDATE "ViewSum1"
  SET
    "Table1Num1"=c111,
    "Table2Num3"=c123,
    "Table1I"=i,
    "Table2I"=i
  WHERE "ID"='1:1';
UPDATE "ViewSum1"
  SET
    "Table1Num1"=c211,
    "Table2Num3"=c223,
    "Table1I"=i,
    "Table2I"=i
  WHERE "ID"='2:2';
END;
$$;
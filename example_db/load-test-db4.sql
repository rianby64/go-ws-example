
  -- This is the function I want to test
DO
$$
DECLARE
  i bigint=0;
BEGIN
  WHILE i < 10000 LOOP
    i = i + 1;
    PERFORM goahead(i);
  END LOOP;
END;
$$;

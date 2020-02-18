DO $$
DECLARE
    r RECORD; 
    cnt INT = 0; 
BEGIN
    FOR r IN 
        SELECT * FROM dept 
    LOOP 
		BEGIN
			if (cnt > 100000) then
				return;
			end if;
			update dept set loc=loc where deptno = r.deptno;
			--commit;
			cnt := cnt + 1;
		EXCEPTION when OTHERS THEN
			RAISE NOTICE 'Illegal operation: %', SQLERRM;
			RETURN;
		END;
	END LOOP; 
END $$;
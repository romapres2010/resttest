CREATE OR REPLACE FUNCTION dept_notify_event() RETURNS TRIGGER AS $$

--    DECLARE 
--        data json;
--        notification json;
    
    BEGIN
    
        -- Convert the old or new row to JSON, based on the kind of action.
        -- Action = DELETE?             -> OLD row
        -- Action = INSERT or UPDATE?   -> NEW row
/*
		IF (TG_OP = 'DELETE') THEN
          data = row_to_json(OLD);
        ELSE
          data = row_to_json(NEW);
        END IF;
*/
/*
	IF (TG_TABLE_NAME= 'dept') THEN
        	IF (TG_OP = 'DELETE') THEN
	           	data = to_json(OLD.detno);
        	ELSE
           		data = to_json(NEW.deptno);
        	END IF;
	END IF;

	IF (TG_TABLE_NAME= 'emp') THEN
        	IF (TG_OP = 'DELETE') THEN
	           	data = to_json(OLD.empno);
        	ELSE
           		data = to_json(NEW.empno);
        	END IF;
	END IF;
	
        -- Contruct the notification as a JSON string.
        notification = json_build_object(
                          'table',TG_TABLE_NAME,
                          'action', TG_OP,
                          'data', data);
*/      
        	IF (TG_OP = 'DELETE') THEN
	           PERFORM pg_notify('events_dept', OLD.deptno::varchar(255));
        	ELSE
	           PERFORM pg_notify('events_dept', NEW.deptno::varchar(255));
        	END IF;
			
        -- Execute pg_notify(channel, notification)
--        PERFORM pg_notify('events_dept', notification::text);
  
        -- Result is ignored since this is an AFTER trigger
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION emp_notify_event() RETURNS TRIGGER AS $$
	DECLARE
    BEGIN
        	IF (TG_OP = 'DELETE') THEN
	           	PERFORM pg_notify('events_emp', OLD.empno::varchar(255));
				PERFORM pg_notify('events_dept', OLD.deptno::varchar(255));
        	ELSEIF (TG_OP = 'INSERT') THEN
	           	PERFORM pg_notify('events_emp', NEW.empno::varchar(255));
			   	PERFORM pg_notify('events_dept', NEW.deptno::varchar(255));
        	ELSEIF (TG_OP = 'UPDATE') THEN
	           	PERFORM pg_notify('events_emp', NEW.empno::varchar(255));	
				-- Если поменялся deptno
				IF NEW.deptno != OLD.deptno THEN
					PERFORM pg_notify('events_dept', OLD.deptno::varchar(255));
			   		PERFORM pg_notify('events_dept', NEW.deptno::varchar(255));
				END IF;
        	END IF;
        RETURN NULL; 
    END;
    
$$ LANGUAGE plpgsql;



CREATE TRIGGER dept_notify_event
AFTER INSERT OR UPDATE OR DELETE ON dept
    FOR EACH ROW EXECUTE PROCEDURE dept_notify_event();

CREATE TRIGGER emp_notify_event
AFTER INSERT OR UPDATE OR DELETE ON emp
    FOR EACH ROW EXECUTE PROCEDURE emp_notify_event();

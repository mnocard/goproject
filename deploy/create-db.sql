DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_database WHERE datname = 'user-db') THEN
      RAISE NOTICE 'Database already exists';  -- optional
   ELSE
      PERFORM dblink_exec('dbname=' || current_database()  -- current db
                        , 'CREATE DATABASE user-db');
   END IF;
END
$do$;

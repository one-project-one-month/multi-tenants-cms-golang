-- 1. First, create the schema and domain correctly
CREATE SCHEMA IF NOT EXISTS Function;

CREATE DOMAIN positive_price AS float CHECK (VALUE >= 0.0);

-- 2. Create the table
CREATE TABLE Function.Courses (
                         course_id serial PRIMARY KEY,
                         course_name varchar(100) NOT NULL,
                         prices positive_price NOT NULL  -- Using the domain
);

-- 3. Proper transaction with error handling
DO $$
    BEGIN
        -- This will succeed
        INSERT INTO Courses (course_name, prices)
        VALUES ('programming', 1000.0);

        -- This will fail but we'll catch the error
        BEGIN
            INSERT INTO Courses (course_name, prices)
            VALUES ('coding', -1000.0);
        EXCEPTION WHEN check_violation THEN
            RAISE NOTICE 'Failed to insert course: negative price not allowed';
        END;

    END $$;

-- 4. Verify the results
SELECT * FROM Courses;
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE  SCHEMA IF NOT EXISTS  Mapping;

CREATE  TABLE  Mapping.Users (
    user_id serial primary key,
    user_name  varchar(100) not null ,
    location  geography(Point, 4326)
);
INSERT INTO Mapping.Users (user_name, location) VALUES
                                       ('User A', ST_SetSRID(ST_MakePoint(-74.0060, 40.7128), 4326)), -- New York
                                       ('User B', ST_SetSRID(ST_MakePoint(2.3522, 48.8566), 4326)),   -- Paris
                                       ('User C', ST_SetSRID(ST_MakePoint(139.6917, 35.6895), 4326)); -- Tokyo

SELECT
    user_name,
    user_name,
    st_x(location::geometry) as longitude,
    st_y(location::geometry) as lattitude
FROM Mapping.Users;
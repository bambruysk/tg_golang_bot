CREATE TABLE IF NOT EXISTS location (id serial PRIMARY KEY, name varchar unique);

CREATE TABLE IF NOT EXISTS players (
    id serial PRIMARY KEY,
    name varchar,
    amount decimal(10, 2),
    last_visit TIMESTAMP
);


CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    name varchar(80) not null,
    location int,

    FOREIGN KEY (location) REFERENCES location (id) ON DELETE
    SET
        NULL
);

CREATE TABLE IF NOT EXISTS holdes (
    id serial PRIMARY KEY,
    name varchar NOT NULL,
    amount decimal,
    level integer not null default 1,
    owner_id integer REFERENCES players (id) ON DELETE SET NULL ,
    last_visit timestamp 
);

-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS songs (
    id serial PRIMARY KEY,
    group_name VARCHAR (1024),
    song VARCHAR (1024),
    release_date date,
    link VARCHAR (1024)
);

CREATE TABLE IF NOT EXISTS verses (
    id serial PRIMARY KEY,
    song_id INT,
    num INT,
    content TEXT
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE songs;
DROP TABLE verses;

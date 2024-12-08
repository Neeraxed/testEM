-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE IF NOT EXISTS songs (
    id int PRIMARY KEY
    group VARCHAR (1024)
    song VARCHAR (1024)
    releaseDate VARCHAR (256)
    link VARCHAR (1024)
);

CREATE TABLE verses (
    id int PRIMARY KEY
    song_id INT
    number INT
    content TEXT
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE songs;
DROP TABLE verses;

-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE INDEX IF NOT EXISTS song_idx on songs using btree (song);
CREATE INDEX IF NOT EXISTS group_idx on songs using btree (group_name);
CREATE INDEX IF NOT EXISTS release_date_idx on songs using btree (release_date);
CREATE INDEX IF NOT EXISTS song_id_idx on verses using btree (song_id);
-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP INDEX IF EXISTS song_idx;
DROP INDEX IF EXISTS group_idx;
DROP INDEX IF EXISTS release_date_idx;
DROP INDEX IF EXISTS song_id_idx;

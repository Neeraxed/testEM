package verse

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (st *Storage) AddVersesForSong(songId string, verses []Verse) error {
	builder := sq.Insert("verses").Columns("song_id", "number", "content")

	for _, verse := range verses {
		builder.Values(songId, verse.Number, verse.Content)
	}

	query, args, err := builder.ToSql()
	_, err = st.db.Exec(query, args...)
	if err != nil {

	}
	return err
}

func (st *Storage) GetVersesFromSong(songId string) ([]Verse, error) {
	builder := sq.Select("*").From("verses").Where(sq.Eq{"song_id": songId})
	queryStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	verses := []Verse{}
	err = st.db.QueryRow(queryStr, args).Scan(&verses)
	return verses, err
}

func (st *Storage) DeleteSong(id string) error {
	builder := sq.Delete("verses").Where(sq.Eq{"song_id": id})
	queryStr, args, err := builder.ToSql()

	if err != nil {
		zap.NamedError("reason", err)
	}

	_, err = st.db.Exec(queryStr, args...)
	if err != nil {

	}
	return err
}

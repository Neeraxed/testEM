package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"testEM/internal/entities"
	"time"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type VerseStorage struct {
	db  *sql.DB
	log *zap.Logger
}

func NewVerseStorage(db *sql.DB, log *zap.Logger) *VerseStorage {
	return &VerseStorage{
		db:  db,
		log: log,
	}
}

func (st *VerseStorage) AddVersesForSong(songId string, verses []*entities.Verse) error {
	builder := sq.Insert("verses").Columns("song_id", "num", "content")

	for _, verse := range verses {
		builder = builder.Values(songId, verse.Number, verse.Content)
	}

	builder = builder.PlaceholderFormat(sq.Dollar)
	query, args, err := builder.ToSql()
	fmt.Println(query)
	if err != nil {
		st.log.Debug("Failed to build sql query to add verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	_, err = st.db.Exec(query, args...)
	if err != nil {
		st.log.Debug("Failed to add text song to verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}

	return err
}

func (st *VerseStorage) GetVersesForSong(opts entities.VerseSearchOptions) ([]*entities.Verse, int, error) {
	builder := sq.Select("*").From("verses")
	builder = st.AddSearchOptionsToBuilder(builder, &opts)
	builder = builder.PlaceholderFormat(sq.Dollar)
	queryStr, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to get verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, 0, err
	}

	var verses []*entities.Verse
	rows, err := st.db.Query(queryStr, args...)

	for rows.Next() {
		v := entities.Verse{}
		var id int
		if err := rows.Scan(&id, &v.SongID, &v.Number, &v.Content); err != nil {
			st.log.Debug("Failed to scan row in GetVersesForSong")
			return nil, 0, err
		}
		verses = append(verses, &v)
	}

	if err = rows.Err(); err != nil {
		st.log.Debug("Failed to scan rows in GetVersesForSong")
		return nil, 0, err
	}
	defer rows.Close()

	builder = sq.Select("count(*)").From("verses")
	builder = st.AddSearchOptionsToBuilder(builder, &opts)
	builder = builder.PlaceholderFormat(sq.Dollar)
	queryStr, args, err = builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to get total verses in GetVersesForSong",
			zap.String("message", err.Error()),
		)
		return nil, 0, err
	}

	var count int
	err = st.db.QueryRow(queryStr, args...).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, &NotFoundErr{}
		}

		st.log.Debug("Failed to execute query fot total verses in GetVersesForSong",
			zap.String("message", err.Error()),
		)
		return nil, 0, err
	}

	return verses, count, err
}

func (st *VerseStorage) DeleteSong(id string) error {
	builder := sq.Delete("verses").Where(sq.Eq{"song_id": id}).PlaceholderFormat(sq.Dollar)
	queryStr, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to delete verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	_, err = st.db.Exec(queryStr, args...)
	if err != nil {
		st.log.Debug("Failed to get text song from verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}
	return err
}

func (st *VerseStorage) AddSearchOptionsToBuilder(builder sq.SelectBuilder, opts *entities.VerseSearchOptions) sq.SelectBuilder {
	if opts.SongID != "" {
		builder = builder.Where(sq.Eq{"song_id": &opts.SongID})
	}
	if opts.Page != 0 && opts.PerPage != 0 {
		builder = builder.Offset((uint64)(opts.PerPage * (opts.Page - 1)))
		builder.Limit(uint64(opts.PerPage))
	}
	return builder
}

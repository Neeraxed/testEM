package repository

import (
	"database/sql"
	"errors"
	"testEM/internal/entities"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type SongStorage struct {
	db  *sql.DB
	log *zap.Logger
}

func NewSongStorage(db *sql.DB, log *zap.Logger) *SongStorage {
	return &SongStorage{
		db:  db,
		log: log,
	}
}

type NotFoundErr struct {
	s string
}

func (e *NotFoundErr) Error() string {
	e.s = "not found"
	return e.s
}

func (st *SongStorage) AddSong(song entities.Song) (*entities.Song, error) {
	builder := sq.Insert("songs").
		Columns("group_name", "song", "release_date", "link").
		Values(*song.Group, *song.Song, *song.ReleaseDate, *song.Link).
		Suffix("RETURNING \"id\"").PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to add song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	err = st.db.QueryRow(query, args...).Scan(&song.ID)
	if err != nil {
		st.log.Debug("Failed to execute query in AddSong",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	return &song, err
}

func (st *SongStorage) GetSongsWithFilters(opts *entities.SongSearchOptions) ([]*entities.Song, int, error) {
	builder := sq.Select("*").From("songs")
	builder = st.AddSearchOptionsToBuilder(builder, opts, true)
	builder = builder.PlaceholderFormat(sq.Dollar)

	queryStr, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to get songs with filters",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, 0, err
	}

	songs := make([]*entities.Song, 0)
	rows, err := st.db.Query(queryStr, args...)

	for rows.Next() {
		s := entities.Song{}
		if err := rows.Scan(&s.ID, &s.Group, &s.Song, &s.ReleaseDate, &s.Link); err != nil {
			st.log.Debug("Failed to scan row in GetSongsWithFilters")
			return nil, 0, err
		}
		songs = append(songs, &s)
	}
	if err = rows.Err(); err != nil {
		st.log.Debug("Failed to scan rows in GetSongsWithFilters")
		return nil, 0, err
	}
	defer rows.Close()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, &NotFoundErr{}
		}

		st.log.Debug("Failed to execute query in GetSongsWithFilters",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, 0, err
	}

	builder = sq.Select("count(*)").From("songs")
	builder = st.AddSearchOptionsToBuilder(builder, opts, false)
	builder = builder.PlaceholderFormat(sq.Dollar)

	queryStr, args, err = builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to get total songs with filters",
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

		st.log.Debug("Failed to execute query for total songs in GetSongsWithFilters",
			zap.String("message", err.Error()),
		)
		return nil, 0, err
	}

	return songs, count, err
}

func (st *SongStorage) DeleteSong(id string) error {
	builder := sq.Delete("songs").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar)

	queryStr, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to delete song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	_, err = st.db.Exec(queryStr, args...)
	if err != nil {
		st.log.Debug("Failed to execute query in DeleteSong",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	return err
}

func (st *SongStorage) UpdateSong(id string, song entities.Song) (*entities.Song, error) {
	builder := sq.Update("songs").Where(sq.Eq{"id": id})
	builder = st.AddUpdateOptionsToBuilder(builder, &song)
	builder = builder.Suffix("RETURNING songs.*").PlaceholderFormat(sq.Dollar)

	queryStr, args, err := builder.ToSql()
	if err != nil {
		st.log.Debug("Failed to build sql query to update song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	err = st.db.QueryRow(queryStr, args...).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Link)
	if err != nil {
		st.log.Debug("Failed to execute query in UpdateSong",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}
	return &song, err
}

func (st *SongStorage) AddSearchOptionsToBuilder(builder sq.SelectBuilder, opts *entities.SongSearchOptions, enablePagination bool) sq.SelectBuilder {
	if opts.Group != nil {
		builder = builder.Where(sq.Eq{"group_name": *opts.Group})
	}

	if opts.Song != nil {
		builder = builder.Where(sq.Eq{"song": *opts.Song})
	}

	if opts.ReleaseDateAfter != nil {
		builder = builder.Where(sq.GtOrEq{"release_date": *opts.ReleaseDateAfter})
	}

	if opts.ReleaseDateBefore != nil {
		builder = builder.Where(sq.LtOrEq{"release_date": *opts.ReleaseDateBefore})
	}

	if enablePagination && opts.Page != nil && opts.PerPage != nil {
		builder = builder.Offset((uint64)(*opts.PerPage * (*opts.Page - 1))).Limit(uint64(*opts.PerPage))
	}
	return builder
}

func (st *SongStorage) AddUpdateOptionsToBuilder(builder sq.UpdateBuilder, song *entities.Song) sq.UpdateBuilder {
	if song.Group != nil {
		builder = builder.Set("group_name", *song.Group)
	}

	if song.Song != nil {
		builder = builder.Set("song", *song.Song)
	}

	if song.Link != nil {
		builder = builder.Set("link", *song.Link)
	}
	if song.ReleaseDate != nil {
		builder = builder.Set("release_date", *song.ReleaseDate)
	}
	return builder
}

package song

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/lib/pq"
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

func (st *Storage) AddSong(song Song) error {
	builder := sq.Insert("songs").Columns("id", "group", "song", "releaseDate", "link").Values(song.ID, song.Group, song.Song, song.ReleaseDate, song.Link)
	query, args, err := builder.ToSql()
	if err != nil {

	}
	_, err = st.db.Exec(query, args...)
	if err != nil {

	}

	return err
}

func (st *Storage) GetSongsWithFilters(opts *SongSearchOptions) ([]*Song, error) {
	builder := sq.Select("*").From("songs")
	st.AddSearchOptionsToBuilder(&builder, opts)

	queryStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var songs []*Song
	err = st.db.QueryRow(queryStr, args).Scan(&songs)
	return songs, err
}

func (st *Storage) DeleteSong(id string) error {
	builder := sq.Delete("songs").Where(sq.Eq{"id": id})
	queryStr, args, err := builder.ToSql()

	if err != nil {
		zap.NamedError("reason", err)
	}

	_, err = st.db.Exec(queryStr, args...)
	if err != nil {

	}
	return err
}

func (st *Storage) UpdateSong(id string, song Song) error {
	builder := sq.Update("songs").Where(sq.Eq{"id": id})
	st.AddUpdateOptionsToBuilder(&builder, &song)

	queryStr, args, err := builder.ToSql()
	if err != nil {

	}
	_, err = st.db.Exec(queryStr, args...)
	return err
}

func (st *Storage) AddSearchOptionsToBuilder(builder *sq.SelectBuilder, opts *SongSearchOptions) {
	if opts.Group != nil {
		builder.Where(sq.Eq{"group": *opts.Group})
	}

	if opts.Song != nil {
		builder.Where(sq.Eq{"song": *opts.Song})
	}

	if opts.ReleaseDateAfter != nil {
		builder.Where(sq.GtOrEq{"releaseDate": *opts.ReleaseDateAfter})
	}

	if opts.ReleaseDateBefore != nil {
		builder.Where(sq.LtOrEq{"releaseDate": *opts.ReleaseDateBefore})
	}

	if opts.Page != nil && opts.PerPage != nil {
		builder.Offset((uint64)(*opts.PerPage * (*opts.Page - 1)))
	}
}

func (st *Storage) AddUpdateOptionsToBuilder(builder *sq.UpdateBuilder, song *Song) {
	if song.Group != nil {
		builder.Set("group", *song.Group)
	}

	if song.Song != nil {
		builder.Set("song", *song.Song)
	}

	if song.Link != nil {
		builder.Set("link", *song.Link)
	}
	if song.ReleaseDate != nil {
		builder.Set("releaseDate", *song.ReleaseDate)
	}
}

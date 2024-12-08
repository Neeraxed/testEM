package song

import "time"

type AddSongDTO struct {
	Group *string
	Song  *string
}

type PatchSongDTO struct {
	Group       *string
	Song        *string
	ReleaseDate *string `json:"releaseDate"`
	Link        *string `json:"link"`
}

type Song struct {
	ID          *string
	Group       *string
	Song        *string
	ReleaseDate *time.Time `json:"releaseDate"`
	Link        *string    `json:"link"`
}

type SongSearchOptions struct {
	Group             *string
	Song              *string
	ReleaseDateBefore *time.Time
	ReleaseDateAfter  *time.Time
	Page              *int
	PerPage           *int
}

// web api
type SongDetail struct {
	ID          string `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Content     string `json:"txt"`
	Link        string `json:"link"`
}

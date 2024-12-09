package entities

import "time"

type AddSongDTO struct {
	Group *string `json:"group"`
	Song  *string `json:"song"`
}

type PatchSongDTO struct {
	Group       *string `json:"group"`
	Song        *string `json:"song"`
	ReleaseDate *string `json:"releaseDate"`
	Link        *string `json:"link"`
}

type Song struct {
	ID          *string    `json:"id"`
	Group       *string    `json:"group"`
	Song        *string    `json:"song"`
	ReleaseDate *time.Time `json:"releaseDate"`
	Link        *string    `json:"link"`
}

type SongSearchOptions struct {
	Group             *string    `json:"group"`
	Song              *string    `json:"song"`
	ReleaseDateBefore *time.Time `json:"releaseDateBefore"`
	ReleaseDateAfter  *time.Time `json:"releaseDateAfter"`
	Page              *int       `json:"page"`
	PerPage           *int       `json:"perPage"`
}

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Content     string `json:"txt"`
	Link        string `json:"link"`
}

type SongsWrapper struct {
	Songs []*Song `json:"songs"`
	Total int     `json:"total"`
}

package entities

type Verse struct {
	SongID  string `json:"song_id"`
	Number  int    `json:"num"`
	Content string `json:"content"`
}

type VerseSearchOptions struct {
	SongID  *string
	Page    *int
	PerPage *int
}
type VersesWrapper struct {
	Verses []*Verse `json:"verses"`
	Total  int      `json:"total"`
}

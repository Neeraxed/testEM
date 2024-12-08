package usecase

import (
	"strings"
	"testEM/internal/repository/song"
	"testEM/internal/repository/verse"
	"time"

	"go.uber.org/zap"
)

const (
	dateLayout = "02.01.2006"
)

type Usecase struct {
	songRepo  SongRepo
	verseRepo VerseRepo
	log       *zap.Logger
	client    DetailClient
}
type DetailClient interface {
	GetSongDetails(song song.AddSongDTO) (song.SongDetail, error)
}

type SongRepo interface {
	GetSongsWithFilters(opts *song.SongSearchOptions) ([]*song.Song, error)
	DeleteSong(id string) error
	UpdateSong(id string, s song.Song) error
	AddSong(track song.Song) error
}

type VerseRepo interface {
	GetVersesFromSong(songId string) ([]verse.Verse, error)
	AddVersesForSong(songId string, verses []verse.Verse) error
	DeleteSong(id string) error
}

// TODO delete check
var _ VerseRepo = &verse.Storage{}
var _ SongRepo = &song.Storage{}

func NewUsecase(sr SongRepo, vr VerseRepo, log *zap.Logger, client DetailClient) *Usecase {
	return &Usecase{
		songRepo:  sr,
		verseRepo: vr,
		log:       log,
		client:    client,
	}
}

func (uc *Usecase) GetSongsWithFilters(options song.SongSearchOptions) ([]*song.Song, error) {
	s, err := uc.songRepo.GetSongsWithFilters(&options)
	if err != nil {

	}

	return s, err
}

func (uc *Usecase) GetSong(id string) ([]verse.Verse, error) {
	verses, err := uc.verseRepo.GetVersesFromSong(id)
	return verses, err
}

func (uc *Usecase) DeleteSong(id string) error {
	err := uc.songRepo.DeleteSong(id)
	if err != nil {

	}
	err = uc.verseRepo.DeleteSong(id)
	if err != nil {

	}
	return err
}

func (uc *Usecase) PatchSong(id string, dto song.PatchSongDTO) error {
	s := song.Song{
		Group: dto.Group,
		Song:  dto.Song,
		Link:  dto.Link,
	}
	if dto.ReleaseDate != nil {
		date, err := time.Parse(dateLayout, *dto.ReleaseDate)
		if err != nil {

		}

		s.ReleaseDate = &date
	}
	err := uc.songRepo.UpdateSong(id, s)

	if err != nil {

	}

	return err
}

func (uc *Usecase) AddSong(dto song.AddSongDTO) error {
	details, err := uc.client.GetSongDetails(dto)
	if err != nil {
		return err
	}
	date, err := time.Parse(dateLayout, details.ReleaseDate)
	if err != nil {

	}
	track := song.Song{
		Group:       &details.Group,
		Song:        &details.Song,
		ReleaseDate: &date,
		Link:        &details.Link,
	}
	err = uc.songRepo.AddSong(track)
	if err != nil {

	}

	content := details.Content

	contents := strings.Split(content, "\n\n")

	var verses []verse.Verse
	for i, entry := range contents {
		verses = append(verses,
			verse.Verse{
				SongID:  details.ID,
				Number:  i + 1,
				Content: entry,
			})
	}

	err = uc.verseRepo.AddVersesForSong(details.ID, verses)
	if err != nil {

	}

	return err
}

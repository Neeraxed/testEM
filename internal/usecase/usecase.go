package usecase

import (
	"strings"
	"testEM/internal/entities"
	"time"

	"go.uber.org/zap"
)

const (
	DateLayout = "02.01.2006"
)

type Usecase struct {
	songRepo  SongRepo
	verseRepo VerseRepo
	log       *zap.Logger
	client    DetailClient
}
type DetailClient interface {
	GetSongDetails(song entities.AddSongDTO) (*entities.SongDetail, error)
}

type SongRepo interface {
	GetSongsWithFilters(opts *entities.SongSearchOptions) ([]*entities.Song, int, error)
	DeleteSong(id string) error
	UpdateSong(id string, s entities.Song) (*entities.Song, error)
	AddSong(song entities.Song) (*entities.Song, error)
}

type VerseRepo interface {
	GetVersesForSong(opts entities.VerseSearchOptions) ([]*entities.Verse, int, error)
	AddVersesForSong(songId string, verses []*entities.Verse) error
	DeleteSong(id string) error
}

func NewUsecase(sr SongRepo, vr VerseRepo, log *zap.Logger, client DetailClient) *Usecase {
	return &Usecase{
		songRepo:  sr,
		verseRepo: vr,
		log:       log,
		client:    client,
	}
}

func (uc *Usecase) GetSongsWithFilters(options entities.SongSearchOptions) (entities.SongsWrapper, error) {
	s, count, err := uc.songRepo.GetSongsWithFilters(&options)
	if err != nil {
		//if errors.Is(err, &repository.NotFoundErr{}) {
		//	uc.log.Error("Songs not found",
		//		zap.String("message", err.Error()),
		//		zap.Time("time", time.Now()),
		//	)
		//}
		uc.log.Error("failed to get songs with filters",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return entities.SongsWrapper{}, err
	}
	uc.log.Info("Recieved list of songs with filters",
		zap.Time("time", time.Now()),
	)
	resp := entities.SongsWrapper{
		Songs: s,
		Total: count,
	}
	return resp, err
}

func (uc *Usecase) GetVerses(options entities.VerseSearchOptions) (entities.VersesWrapper, error) {
	verses, count, err := uc.verseRepo.GetVersesForSong(options)
	if err != nil {
		uc.log.Error("failed to get verses for song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return entities.VersesWrapper{}, err
	}

	uc.log.Info("Recieved song text",
		zap.Time("time", time.Now()),
	)
	resp := entities.VersesWrapper{
		Verses: verses,
		Total:  count,
	}
	return resp, err
}

func (uc *Usecase) DeleteSong(id string) error {
	err := uc.songRepo.DeleteSong(id)
	if err != nil {
		uc.log.Error("Failed to delete song from songs",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	uc.log.Info("Deleted song from songs",
		zap.Time("time", time.Now()),
	)

	err = uc.verseRepo.DeleteSong(id)
	if err != nil {
		uc.log.Error("Failed to delete song text from verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return err
	}

	uc.log.Info("Deleted song text from verses",
		zap.Time("time", time.Now()),
	)

	return err
}

func (uc *Usecase) PatchSong(id string, dto entities.PatchSongDTO) (*entities.Song, error) {
	s := entities.Song{
		Group: dto.Group,
		Song:  dto.Song,
		Link:  dto.Link,
	}
	if dto.ReleaseDate != nil {
		date, err := time.Parse(DateLayout, *dto.ReleaseDate)
		if err != nil {
			uc.log.Error("Failed to read release date",
				zap.String("message", err.Error()),
				zap.Time("time", time.Now()),
			)
		}
		s.ReleaseDate = &date
	}
	resp, err := uc.songRepo.UpdateSong(id, s)

	if err != nil {
		uc.log.Error("Failed to update song in songs",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	uc.log.Info("Updated song",
		zap.Time("time", time.Now()),
	)

	return resp, err
}

func (uc *Usecase) AddSong(dto entities.AddSongDTO) (*entities.Song, error) {
	details, err := uc.client.GetSongDetails(dto)
	if err != nil {
		uc.log.Error("Failed to get song details from external API",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	uc.log.Info("Recieved song details from external API",
		zap.Time("time", time.Now()),
	)

	date, err := time.Parse(DateLayout, details.ReleaseDate)
	if err != nil {
		uc.log.Error("Failed to parse release date",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}

	track := entities.Song{
		Group:       dto.Group,
		Song:        dto.Song,
		ReleaseDate: &date,
		Link:        &details.Link,
	}

	s, err := uc.songRepo.AddSong(track)
	if err != nil {
		uc.log.Error("Failed to add song in songs",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	uc.log.Info("Added song to songs",
		zap.Time("time", time.Now()),
	)

	content := details.Content
	contents := strings.Split(content, "\n\n")

	var verses []*entities.Verse
	for i, entry := range contents {
		verses = append(verses,
			&entities.Verse{
				SongID:  *s.ID,
				Number:  i + 1,
				Content: entry,
			})
	}

	err = uc.verseRepo.AddVersesForSong(*s.ID, verses)
	if err != nil {
		uc.log.Error("Failed to add song text in verses",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)

		return nil, err
	}

	uc.log.Info("Added song text to verses",
		zap.Time("time", time.Now()),
	)

	return s, err
}

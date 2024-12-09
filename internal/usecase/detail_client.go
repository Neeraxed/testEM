package usecase

import (
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"testEM/internal/entities"
	"time"
)

type detailClient struct {
	externalUrl string
	client      *http.Client
	log         *zap.Logger
}

func NewDetailClient(externalUrl string, client *http.Client, log *zap.Logger) DetailClient {
	return &detailClient{
		externalUrl: externalUrl,
		client:      client,
		log:         log,
	}
}

func (dt *detailClient) GetSongDetails(track entities.AddSongDTO) (*entities.SongDetail, error) {
	u, err := url.Parse(dt.externalUrl)
	if err != nil {
		dt.log.Debug("Failed to read external url",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	u.Query().Add("song", *track.Song)
	u.Query().Add("group", *track.Group)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		dt.log.Debug("Failed to create request",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	resp, err := dt.client.Do(req)
	if err != nil {
		dt.log.Debug("Failed to execute request",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		dt.log.Debug("Failed to read response body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	var detail *entities.SongDetail

	err = json.Unmarshal(body, &detail)
	if err != nil {
		dt.log.Debug("Failed to unmarshal response body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return nil, err
	}

	return detail, err
}

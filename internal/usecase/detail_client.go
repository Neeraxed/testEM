package usecase

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testEM/internal/repository/song"
)

type detailClient struct {
	externalUrl string
	client      *http.Client
}

func (dt *detailClient) GetSongDetails(track song.AddSongDTO) (*song.SongDetail, error) {
	u, err := url.Parse(dt.externalUrl)
	if err != nil {
		return nil, err
	}
	u.Query().Add("song", *track.Song)
	u.Query().Add("group", *track.Group)

	req, err := http.NewRequest("GET", u.String(), nil)
	resp, err := dt.client.Do(req)

	body, err := io.ReadAll(resp.Body)
	{
		if err != nil {
			return nil, err
		}
	}

	var detail *song.SongDetail

	err = json.Unmarshal(body, &detail)

	return detail, err
}

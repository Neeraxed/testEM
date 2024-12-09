package delivery

import "testEM/internal/entities"

type MockExternal struct {
}

func (m *MockExternal) GetSongDetails(track entities.AddSongDTO) (*entities.SongDetail, error) {
	return &entities.SongDetail{
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
		ReleaseDate: "16.07.2006",
		Content:     "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight\n",
	}, nil
}

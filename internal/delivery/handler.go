package delivery

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"testEM/internal/repository/song"
	"testEM/internal/usecase"
	"testEM/pkg/middleware"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	songsUrl  = "api/v1/songs"
	songUrl   = "api/v1/songs/:id"
	versesUrl = "api/v1/verses/:id"
)

type Handler interface {
	ApplyRoutes(o *middleware.Onion) *httprouter.Router
}

type handler struct {
	Logger *zap.Logger
	uc     *usecase.Usecase
}

func NewHandler(lg *zap.Logger, uc *usecase.Usecase) Handler {
	return &handler{
		Logger: lg,
		uc:     uc,
	}
}

func (h *handler) ApplyRoutes(o *middleware.Onion) *httprouter.Router {
	router := httprouter.New()
	router.GET(songsUrl, o.Apply(h.GetSongs))
	router.GET(versesUrl, o.Apply(h.GetTextBySongID))
	router.DELETE(songUrl, o.Apply(h.DeleteSong))
	router.PATCH(songUrl, o.Apply(h.PatchSong))
	router.POST(songUrl, o.Apply(h.AddSong))

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return router
}

func (h *handler) GetSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	searchOptions := song.SongSearchOptions{}

	myUrl, err := url.Parse(r.URL.Path)
	if err != nil {

	}
	params := myUrl.Query()

	if params != nil && len(params) > 0 {
		if val := params.Get("song"); val != "" {
			searchOptions.Song = &val
		}
		if val := params.Get("group"); val != "" {
			searchOptions.Group = &val
		}
		if val := params.Get("page"); val != "" {
			page, err := strconv.Atoi(val)
			if err != nil {

			}
			searchOptions.Page = &page
		}
		if val := params.Get("perpage"); val != "" {
			perpage, err := strconv.Atoi(val)
			if err != nil {

			}
			searchOptions.PerPage = &perpage
		}

		//TODO add date filter
	}

	s, err := h.uc.GetSongsWithFilters(searchOptions)
	if err != nil {

	}
	resp, err := json.Marshal(s)
	if err != nil {

	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) GetTextBySongID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	songID := ps.ByName("id")

	s, err := h.uc.GetSong(songID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	resp, err := json.Marshal(s)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) DeleteSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	songID := ps.ByName("id")
	err := h.uc.DeleteSong(songID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *handler) PatchSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	songID := ps.ByName("id")
	patchDTO := song.PatchSongDTO{}

	myUrl, err := url.Parse(r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	params := myUrl.Query()

	if params != nil && len(params) > 0 {
		if val := params.Get("song"); val != "" {
			patchDTO.Song = &val
		}
		if val := params.Get("group"); val != "" {
			patchDTO.Group = &val
		}
		if val := params.Get("link"); val != "" {
			patchDTO.Link = &val
		}
		if val := params.Get("releaseDate"); val != "" {
			patchDTO.ReleaseDate = &val
		}
	}

	err = h.uc.PatchSong(songID, patchDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *handler) AddSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	s := ps.ByName("id")
	group := ps.ByName("group")
	songDTO := song.AddSongDTO{
		Group: &group,
		Song:  &s,
	}
	err := h.uc.AddSong(songDTO)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

package delivery

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testEM/internal/entities"
	"testEM/internal/repository"
	"testEM/internal/usecase"
	"testEM/pkg/middleware"
	"time"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

const (
	songsUrl  = "/api/v1/songs"
	songUrl   = "/api/v1/songs/:id"
	versesUrl = "/api/v1/songs/:id/verses"
)

type Handler interface {
	ApplyRoutes(o *middleware.Onion) *httprouter.Router
}

type handler struct {
	log *zap.Logger
	uc  *usecase.Usecase
}

func NewHandler(lg *zap.Logger, uc *usecase.Usecase) Handler {
	return &handler{
		log: lg,
		uc:  uc,
	}
}

func (h *handler) ApplyRoutes(o *middleware.Onion) *httprouter.Router {
	router := httprouter.New()
	router.GET(songsUrl, o.Apply(h.GetSongs))
	router.GET(versesUrl, o.Apply(h.GetVersesBySongID))
	router.DELETE(songUrl, o.Apply(h.DeleteSong))
	router.PATCH(songUrl, o.Apply(h.PatchSong))
	router.POST(songsUrl, o.Apply(h.AddSong))

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(http.StatusNoContent)
		return
	})

	return router
}

func (h *handler) GetSongs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	searchOptions := entities.SongSearchOptions{}

	params := r.URL.Query()
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
				h.log.Debug("Failed to parse page int from string",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.Page = &page
		}
		if val := params.Get("perPage"); val != "" {
			perpage, err := strconv.Atoi(val)
			if err != nil {
				h.log.Debug("Failed to parse page int from string",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.PerPage = &perpage
		}

		if val := params.Get("releaseDateAfter"); val != "" {
			d, err := time.Parse(usecase.DateLayout, val)
			if err != nil {
				h.log.Debug("Failed to parse release date",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.ReleaseDateAfter = &d
		}
		if val := params.Get("releaseDateBefore"); val != "" {
			d, err := time.Parse(usecase.DateLayout, val)
			if err != nil {
				h.log.Debug("Failed to parse release date",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.ReleaseDateBefore = &d
		}
	}

	s, err := h.uc.GetSongsWithFilters(searchOptions)
	if err != nil {
		if errors.Is(err, &repository.NotFoundErr{}) {
			h.log.Error("Failed get songs with filters: not found",
				zap.String("message", err.Error()),
				zap.Time("time", time.Now()),
			)
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			h.log.Error("Failed get songs with filters",
				zap.String("message", err.Error()),
				zap.Time("time", time.Now()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp, err := json.Marshal(s)
	if err != nil {
		h.log.Debug("Failed to serialize response",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) GetVersesBySongID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	songID := ps.ByName("id")
	searchOptions := entities.VerseSearchOptions{}
	searchOptions.SongID = songID

	params := r.URL.Query()
	if params != nil && len(params) > 0 {
		if val := params.Get("page"); val != "" {
			page, err := strconv.Atoi(val)
			if err != nil {
				h.log.Debug("Failed to parse page int from string",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.Page = page
		}
		if val := params.Get("perPage"); val != "" {
			perpage, err := strconv.Atoi(val)
			if err != nil {
				h.log.Debug("Failed to parse page int from string",
					zap.String("message", err.Error()),
					zap.Time("time", time.Now()),
				)
			}
			searchOptions.PerPage = perpage
		}
	}

	s, err := h.uc.GetVerses(searchOptions)
	if err != nil {
		h.log.Error("Failed get song text",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp, err := json.Marshal(s)
	if err != nil {
		h.log.Debug("Failed to serialize response",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) DeleteSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	songID := ps.ByName("id")
	err := h.uc.DeleteSong(songID)
	if err != nil {
		h.log.Error("Failed delete song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) PatchSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	songID := ps.ByName("id")
	patchDTO := entities.PatchSongDTO{}

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &patchDTO)
	if err != nil {
		h.log.Error("Failed to read body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	song, err := h.uc.PatchSong(songID, patchDTO)
	if err != nil {
		h.log.Error("Failed to update song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	resp, err := json.Marshal(song)
	if err != nil {
		h.log.Debug("Failed to serialize response",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

func (h *handler) AddSong(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	songDTO := entities.AddSongDTO{}
	err = json.Unmarshal(body, &songDTO)
	if err != nil {
		h.log.Error("Failed to read body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	song, err := h.uc.AddSong(songDTO)
	if err != nil {
		h.log.Error("Failed to add song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	resp, err := json.Marshal(song)
	if err != nil {
		h.log.Debug("Failed to serialize response",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}

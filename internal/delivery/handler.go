package delivery

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi"
	"io"
	"net/http"
	"strconv"
	"testEM/internal/entities"
	"testEM/internal/repository"
	"testEM/internal/usecase"
	"testEM/pkg/middleware"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
	_ "testEM/docs"
)

const (
	songsUrl  = "/api/v1/songs"
	songUrl   = "/api/v1/songs/{id}"
	versesUrl = "/api/v1/songs/{id}/verses"
)

type Handler interface {
	ApplyRoutes(o *middleware.Onion) *chi.Mux
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

func (h *handler) ApplyRoutes(o *middleware.Onion) *chi.Mux {
	router := chi.NewRouter()
	router.Get(songsUrl, o.Apply(h.GetSongs))
	router.Get(versesUrl, o.Apply(h.GetVersesBySongID))
	router.Delete(songUrl, o.Apply(h.DeleteSong))
	router.Patch(songUrl, o.Apply(h.PatchSong))
	router.Post(songsUrl, o.Apply(h.AddSong))

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))
	return router
}

// @Summary      Get songs
// @Description  get string by filters
// @Produce      json
// @Success      200  {object}  entities.SongsWrapper
// @Failure      400  {object} HttpError
// @Failure      404  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /songs [get]
func (h *handler) GetSongs(w http.ResponseWriter, r *http.Request) {
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
			ReturnHttpError(w, err)
			return
		} else {
			h.log.Error("Failed get songs with filters",
				zap.String("message", err.Error()),
				zap.Time("time", time.Now()),
			)
			w.WriteHeader(http.StatusInternalServerError)
			ReturnHttpError(w, err)
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
		ReturnHttpError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// @Summary      Get verses
// @Description  get verses for song
// @Produce      json
// @Success      200  {object}  entities.VersesWrapper
// @Failure      400  {object} HttpError
// @Failure      404  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /songs/{id}/verses [get]
func (h *handler) GetVersesBySongID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	songID := chi.URLParam(r, "id")
	searchOptions := entities.VerseSearchOptions{}
	searchOptions.SongID = &songID

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
	}

	s, err := h.uc.GetVerses(searchOptions)
	if err != nil {
		h.log.Error("Failed get song text",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusNotFound)
		ReturnHttpError(w, err)
		return
	}
	resp, err := json.Marshal(s)
	if err != nil {
		h.log.Debug("Failed to serialize response",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		ReturnHttpError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// @Summary      Delete song
// @Description  delete song with specified id
// @Produce      json
// @Success      200  {object} nil
// @Failure      400  {object} HttpError
// @Failure      404  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /songs [delete]
func (h *handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	songID := chi.URLParam(r, "id")
	err := h.uc.DeleteSong(songID)
	if err != nil {
		h.log.Error("Failed delete song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		ReturnHttpError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// @Summary      Patch song
// @Description  update song with specified id
// @Produce      json
// @Success      200  {object} entities.Song
// @Failure      400  {object} HttpError
// @Failure      404  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /songs [patch]
func (h *handler) PatchSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	songID := chi.URLParam(r, "id")
	patchDTO := entities.PatchSongDTO{}

	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &patchDTO)
	if err != nil {
		h.log.Error("Failed to read body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusBadRequest)
		ReturnHttpError(w, err)
		return
	}

	song, err := h.uc.PatchSong(songID, patchDTO)
	if err != nil {
		h.log.Error("Failed to update song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusNotFound)
		ReturnHttpError(w, err)
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
		ReturnHttpError(w, err)
		return
	}
	w.Write(resp)
}

// @Summary      Add song
// @Description  add song
// @Produce      json
// @Success      200  {object} entities.Song
// @Failure      400  {object} HttpError
// @Failure      404  {object} HttpError
// @Failure      500  {object} HttpError
// @Router       /songs [post]
func (h *handler) AddSong(w http.ResponseWriter, r *http.Request) {
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
		ReturnHttpError(w, err)
		return
	}

	song, err := h.uc.AddSong(songDTO)
	if err != nil {
		h.log.Error("Failed to add song",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		w.WriteHeader(http.StatusInternalServerError)
		ReturnHttpError(w, err)
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
		ReturnHttpError(w, err)
		return
	}
	w.Write(resp)
}

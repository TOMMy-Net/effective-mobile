package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/TOMMy-Net/effective-mobile/internal/storage/db"
	"github.com/TOMMy-Net/effective-mobile/models"
	"github.com/TOMMy-Net/effective-mobile/tools"
	"github.com/TOMMy-Net/effective-mobile/tools/verse"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Service struct {
	Storage *db.Storage
	Log     *logrus.Logger
}

var (
	ErrValid      = errors.New("required fields are not filled in")
	ErrJSON       = errors.New("not volid data")
	ErrSetData    = errors.New("impossible to enter data")
	ErrDeleteData = errors.New("data cannot be deleted")
	ErrBadMethod  = errors.New("method not allowed")
	ErrTypeID     = errors.New("id type not correct")
	ErrData       = errors.New("data not correct, (0000-00-00)")
	ErrWithQuery  = errors.New("error with data in query")
	ErrGetData    = errors.New("data cannot be get")
	ErrNoID       = errors.New("no such id in the database")
)

const (
	DeleteOK = "data be deleted"
	EditOK   = "data be edited"
)

func (s *Service) SongHandlers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			s.GetFilterSongHandler().ServeHTTP(w, r)
		case http.MethodPost:
			s.AddSongHandler().ServeHTTP(w, r)
		case http.MethodDelete:
			s.DeleteSongHandler().ServeHTTP(w, r)
		case http.MethodPatch:
			s.EditSongHandler().ServeHTTP(w, r)
		}
	}
}

// @Summary Add song
// @Tags add
// @Description add song to library
// @ID add-song
// @Accept json
// @Produce json
// @Param input body models.Song true "song info"
// @Success 200 {object} models.Song
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Failure 502 {object} tools.Error
// @Router /api/v1/songs [post]
func (s *Service) AddSongHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var song models.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {

			tools.SetJSON(400, tools.Error{
				Msg: ErrJSON.Error(),
			}, w)
			return
		}

		if err := tools.ValidateFields(&song, "Song", "Group", "ReleaseDate"); err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrValid.Error(),
			}, w)
			return
		}

		t, err := time.Parse(time.DateOnly, song.ReleaseDate)
		if err == nil {
			song.ReleaseDate = t.Format("2006-01-02")
		} else {
			tools.SetJSON(400, tools.Error{
				Msg: ErrData.Error(),
			}, w)
			return
		}

		err = s.Storage.AddSong(r.Context(), &song)
		if err != nil {
			tools.SetJSON(500, tools.Error{
				Msg: ErrSetData.Error(),
			}, w)
			return
		}

		api, err := tools.GetMusicInfo(song.Song, song.Group)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"song":  song.Song,
				"group": song.Group,
				"error": err,
			}).Debug("api error")

			tools.SetJSON(502, tools.Error{
				Msg: err.Error(),
			}, w)
			return
		}
		tools.SetJSON(200, api, w)
	}
}

// @Summary Delete song
// @Tags delete
// @Description delete song from library
// @ID delete-song
// @Accept json
// @Produce json
// @Param id query int true "song id"
// @Success 200 {object} tools.OK
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Router /api/v1/songs [delete]
func (s *Service) DeleteSongHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id = r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrTypeID.Error(),
			}, w)
			return
		}

		if err := s.Storage.DeleteSong(r.Context(), idInt); err != nil {
			tools.SetJSON(500, tools.Error{
				Msg: ErrDeleteData.Error(),
			}, w)
			return
		}

		tools.SetJSON(200, tools.OK{
			Msg: DeleteOK,
		}, w)
	}
}

// @Summary Edit song
// @Tags edit
// @Description edit song at library
// @ID edit-song
// @Accept json
// @Produce json
// @Param input body models.Song true "song info"
// @Param id query int true "song id"
// @Success 200 {object} tools.OK
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Router /api/v1/songs [patch]
func (s *Service) EditSongHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var song models.Song
		var id = r.URL.Query().Get("id")
		idInt, err := strconv.Atoi(id)
		if err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrTypeID.Error(),
			}, w)
			return
		}
		song.ID = idInt

		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrJSON.Error(),
			}, w)
			return
		}

		if err := tools.ValidateFields(&song, "Song", "Group", "ReleaseDate"); err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrValid.Error(),
			}, w)
			return
		}

		err = s.Storage.EditSong(r.Context(), &song)
		if err != nil {
			tools.SetJSON(500, tools.Error{
				Msg: ErrSetData.Error(),
			}, w)
			return
		}
		tools.SetJSON(200, tools.OK{Msg: EditOK}, w)

	}
}

// @Summary Get verse text of song
// @Tags get
// @Description get verse text of song from the library
// @ID getText-song
// @Accept json
// @Produce json
// @Param verse query int true "verse pagination"
// @Param id path int true "song id"
// @Success 200 {object} models.Song
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Router /api/v1/songs/{id}/text [get]
func (s *Service) GetSongTextHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		var page int // куплет песни

		id, errID := strconv.Atoi(mux.Vars(r)["id"])
		page, errPage := strconv.Atoi(r.URL.Query().Get("verse"))
		if errID != nil || errPage != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrWithQuery.Error(),
			}, w)
			return
		}

		text, err := s.Storage.GetSongText(r.Context(), id)
		if err != nil {
			if err == sql.ErrNoRows {
				tools.SetJSON(400, tools.Error{
					Msg: ErrNoID.Error(),
				}, w)
				return
			} else {
				tools.SetJSON(500, tools.Error{
					Msg: ErrGetData.Error(),
				}, w)
				return
			}
		}

		textWithVerse, err := verse.TextPaginate(text, page)
		if err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: verse.ErrNoVerse.Error(),
			}, w)
			return
		}

		tools.SetJSON(200, models.Song{Text: textWithVerse, ID: id}, w)
	}
}

func NewService() *Service {
	return &Service{}
}

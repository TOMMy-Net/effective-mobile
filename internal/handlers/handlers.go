package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/TOMMy-Net/effective-mobile/internal/storage/db"
	"github.com/TOMMy-Net/effective-mobile/models"
	"github.com/TOMMy-Net/effective-mobile/tools"
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
// @Router /songs [post]
func (s *Service) AddSongHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var song models.Song
		if err := json.NewDecoder(r.Body).Decode(&song); err != nil {

			tools.SetJSON(400, tools.Error{
				Msg: ErrJSON.Error(),
			}, w)
			return
		}

		if err := tools.ValidateFields(&song, "Song", "Group"); err != nil {
			tools.SetJSON(400, tools.Error{
				Msg: ErrValid.Error(),
			}, w)
			return
		}

		if song.ReleaseDate != "" {
			t, err := time.Parse(time.DateOnly, song.ReleaseDate)
			if err == nil {
				song.ReleaseDate = t.Format("2006-01-02")
			}
		}

		err := s.Storage.AddSong(r.Context(), &song)
		if err != nil {
			tools.SetJSON(500, tools.Error{
				Msg: ErrSetData.Error(),
			}, w)
			return
		}

		api, err := tools.GetMusicInfo(song.Song, song.Group)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"song": song.Song,
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
// @Tags Delete
// @Description delete song from library
// @ID delete-song
// @Accept json
// @Produce json
// @Param id query int true "song id"
// @Success 200 {object} tools.OK
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Router /songs [delete]
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
// @Success 200 {object} tools.OK
// @Failure 400 {object} tools.Error
// @Failure 500 {object} tools.Error
// @Router /songs [patch]
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

func NewService() *Service {
	return &Service{}
}

package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/TOMMy-Net/effective-mobile/tools"
	"github.com/TOMMy-Net/effective-mobile/tools/filters"
)

var (
	ErrNotPag          = errors.New("pagination is required")
	ErrGetFilteredData = errors.New("error with get data")
)
var (
	pageSize = 10
)

func (s *Service) GetFilterSongHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter := []filters.Filter{}
		var pag int

		var id = r.URL.Query().Get("id")
		var song = r.URL.Query().Get("song")
		var group = r.URL.Query().Get("music_group")
		var link = r.URL.Query().Get("link")
		var text = r.URL.Query().Get("text")
		var date = r.URL.Query().Get("release_date")
		var page = r.URL.Query().Get("page")

		if page != "" {
			page, err := strconv.Atoi(page)
			if err == nil {
				pag = page
			} else {
				tools.SetJSON(400, tools.Error{
					Msg: ErrNotPag.Error(),
				}, w)
				return
			}
		} else {
			tools.SetJSON(400, tools.Error{
				Msg: ErrNotPag.Error(),
			}, w)
			return
		}

		if id != "" {
			operator, value := filters.GetOperatorAndValue(id)
			if value, err := strconv.Atoi(value); err == nil {
				filter = append(filter, filters.Filter{
					Operator: operator,
					Value:    value,
					Name:     "id",
				})
			}
		}

		if song != "" {
			operator, value := filters.GetOperatorAndValue(song)
			filter = append(filter, filters.Filter{
				Operator: operator,
				Value:    value,
				Name:     "song",
			})
		}

		if group != "" {
			operator, value := filters.GetOperatorAndValue(group)
			filter = append(filter, filters.Filter{
				Operator: operator,
				Value:    value,
				Name:     "music_group",
			})
		}

		if link != "" {
			filter = append(filter, filters.Filter{
				Operator: filters.OperatorEq,
				Value:    link,
				Name:     "link",
			})
		}

		if text != "" {
			filter = append(filter, filters.Filter{
				Operator: filters.OperatorEq,
				Value:    text,
				Name:     "text",
			})
		}

		if date != "" {
			operator, value := filters.GetOperatorAndValue(date)
			time, err := time.Parse(time.DateOnly, value)
			if err == nil {
				filter = append(filter, filters.Filter{
					Operator: operator,
					Value:    time.Format("2006-01-02"),
					Name:     "release_date",
				})
			}
		}

		settings := filters.FilterSettings{
			F:          filter,
			Pagination: pag,
			PageSize:   pageSize,
			FieldOrder: "id",
		}

		query := settings.GetFilterWithPagination()

		songs, err := s.Storage.SelectMusicWithCondition(r.Context(), query)
		if err != nil {
			tools.SetJSON(500, tools.Error{
				Msg: ErrGetFilteredData.Error(),
			}, w)
			return
		}

		tools.SetJSON(200, &songs, w)
	}
}

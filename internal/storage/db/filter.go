package db

import (
	"context"
	"fmt"

	"github.com/TOMMy-Net/effective-mobile/models"
)

func (s *Storage) SelectMusicWithCondition(ctx context.Context, where string) ([]models.Song, error) {
	var songs []models.Song
	query := fmt.Sprint("SELECT id, song, music_group, text, link, TO_CHAR(release_date, 'YYYY-MM-DD')  FROM musics " + where)
	p, err := s.DB.PreparexContext(ctx, query)
	if err != nil {
		return []models.Song{}, err
	}

	row, err := p.QueryxContext(ctx)
	if err != nil {
		return []models.Song{}, err
	}

	for row.Next() {
		var song models.Song
		err = row.Scan(&song.ID, &song.Song, &song.Group, &song.Text, &song.Link, &song.ReleaseDate)
		if err == nil {
			songs = append(songs, song)
		}
	}

	return songs, err
}

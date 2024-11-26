package db

import (
	"context"

	"github.com/TOMMy-Net/effective-mobile/models"
)

func (s *Storage) AddSong(ctx context.Context, m *models.Song) error {
	_, err := s.DB.NamedExecContext(ctx, `INSERT INTO musics (song, music_group, text, link, release_date) VALUES (:song, :music_group, :text, :link, :release_date);`, m)
	return err
}

func (s *Storage) DeleteSong(ctx context.Context, id int) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM musics WHERE id = $1`, id)
	return err
}

func (s *Storage) EditSong(ctx context.Context, m *models.Song) error {
	_, err := s.DB.NamedExecContext(ctx, `UPDATE musics SET
												song = :song, 
												music_group = :music_group, 
												text = :text, 
												link = :link, 
												release_date = :release_date WHERE id = :id`, m)
										
	return err
}

package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/TOMMy-Net/effective-mobile/models"
)

var (
	ErrBadReq = errors.New("bad request")
)

func GetMusicInfo(song, group string) (*models.Song, error) {
	resp, err := http.Get(fmt.Sprintf("%s/info?group=%s&song=%s", os.Getenv("EXTERNAL_API"), group, song))
	if err != nil {
		return &models.Song{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return &models.Song{}, ErrBadReq
		} else {
			return &models.Song{}, fmt.Errorf("status code (%d)", resp.StatusCode)
		}
	}
	
	var songInfo models.Song
	if err = json.NewDecoder(resp.Body).Decode(&songInfo); err != nil {
		return &models.Song{}, err
	}
	return &songInfo, nil
}

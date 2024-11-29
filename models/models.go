package models

// for json
type Song struct {
	ID          int    `json:"id" db:"id"`
	Song        string `json:"song" db:"song" validate:"required,min=1"`
	Group       string `json:"group" db:"music_group" validate:"required,min=1"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
	ReleaseDate string `json:"releaseDate" db:"release_date"`
}

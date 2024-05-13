package model

import "time"

type Song struct {
	ID          int64
	Path        string
	LastUpdate  time.Time
	Title       string
	Artist      string
	Album       string
	Genre       string
	ReleaseYear string
}

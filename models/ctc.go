package models

import "time"

type Status string

const (
	Pending Status = "pending"
	Ready   Status = "ready"
)

type CTC struct {
	LongURL     string    `json:"long_url"`
	ReleaseDate time.Time `json:"release_date"`
	Shortcode   string    `json:"shortcode"`
	Status      Status    `json:"status"`
}

var CTCStore = make(map[string]CTC)

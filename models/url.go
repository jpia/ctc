package models

import "time"

type Status string

const (
	PendingStatus  Status = "pending"
	DelayedStatus  Status = "delayed"
	ReleasedStatus Status = "released"
)

type ReleaseMethod string

const (
	OverrideReleaseMethod   ReleaseMethod = "override"
	StandardReleaseMethod   ReleaseMethod = "standard"
	ApiSickDayReleaseMethod ReleaseMethod = "api_sick_day"
)

type URL struct {
	LongURL          string        `json:"long_url"`
	ReleaseDate      time.Time     `json:"release_date"`
	ReleaseDateOrig  time.Time     `json:"release_date_orig"`
	Shortcode        string        `json:"shortcode"`
	Status           Status        `json:"status"`
	Delays           int           `json:"delays"`
	ReleaseMethod    ReleaseMethod `json:"release_method"`
	ReleaseTimestamp time.Time     `json:"release_timestamp"`
}

var URLStore = make(map[string]URL)

package model

import "time"

type Link struct {
	ID          int64
	ShortCode   string
	OriginalURL string
	CreatedAt   time.Time
	Visits      int64
}

package models

import "time"

type PRStatus string

const (
	OpenStatus   PRStatus = "OPEN"
	MergedStatus PRStatus = "MERGED"
)

type PRShort struct {
	PRID     string
	PRName   string
	AuthorID string
	Status   PRStatus
}

type PR struct {
	PRID              string
	PRName            string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          time.Time
}

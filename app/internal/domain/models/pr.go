package models

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

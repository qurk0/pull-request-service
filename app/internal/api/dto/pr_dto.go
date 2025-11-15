package dto

import "time"

type PRShort struct {
	ID       string `json:"pull_request_id"`
	Name     string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
	Status   string `json:"status"`
}

type CreatePRRequest struct {
	PRID     string `json:"pull_request_id"`
	PRNamme  string `json:"pull_request_name"`
	AuthorID string `json:"author_id"`
}

type ReassignPRRequest struct {
	PRID          string `json:"pull_request_id"`
	OldReviewerID string `json:"old_user_id"`
}

type ReassignPRResponse struct {
	Pr         PR     `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}

type PRMergeRequest struct {
	PRID string `json:"pull_request_id"`
}

type PR struct {
	PRID              string     `json:"pull_request_id"`
	PRName            string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt,omitempty"`
}

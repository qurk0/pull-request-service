package dto

type SetIsActiveRequest struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetIsActiveResponse struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type GetReviewResponse struct {
	UserID        string    `json:"user_id"`
	RequestsShort []PRShort `json:"pull_requests"`
}

package dto

type GetTeamResponse struct {
	TeamName string    `json:"team_name"`
	Members  []*Member `json:"members"`
}

type Member struct {
	Id       string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

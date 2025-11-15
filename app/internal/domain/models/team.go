package models

type Team struct {
	TeamName    string
	TeamMembers []TeamMember
}

type TeamMember struct {
	Id       string
	Username string
	IsActive bool
}

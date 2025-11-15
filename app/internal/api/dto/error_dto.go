package dto

/*
   Ошибки документации:
   - TEAM_EXISTS
   - PR_EXISTS
   - PR_MERGED
   - NOT_ASSIGNED
   - NO_CANDIDATE
   - NOT_FOUND
   - UNAUTHORIZED
*/

type ErrCode string

const (
	ErrCodeTeamExists  ErrCode = "TEAM_EXISTS"
	ErrCodePrExists    ErrCode = "PR_EXISTS"
	ErrCodePrMerged    ErrCode = "PR_MERGED"
	ErrCodeNotAssigned ErrCode = "NOT_ASSIGNED"
	ErrCodeNoCandidate ErrCode = "NO_CANDIDATE"
	ErrCodeNotFound    ErrCode = "NOT_FOUND"
)

type ErrorResponse struct {
	Error HttpError `json:"error"`
}

type HttpError struct {
	Code    ErrCode `json:"code"`
	Message string  `json:"message"`
}

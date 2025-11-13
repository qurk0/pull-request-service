package dto

type ErrorResponse struct {
	Error AppError `json:"error"`
}

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

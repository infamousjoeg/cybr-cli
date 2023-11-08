package shared

// ErrorResponse from Logon
type ErrorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

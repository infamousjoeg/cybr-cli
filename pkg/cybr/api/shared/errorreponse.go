package shared

//ErrorReponse from logon
type ErrorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

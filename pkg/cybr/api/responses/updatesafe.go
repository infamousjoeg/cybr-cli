package responses

// UpdateSafe contains the response to the Update Safe function's request
type UpdateSafe struct {
	SafeName                  string `json:"SafeName"`
	Description               string `json:"Description"`
	NumberOfDaysRetention     int    `json:"NumberOfDaysRetention"`
	NumberOfVersionsRetention int    `json:"NumberOfVersionsRetention"`
	OLACEnabled               bool   `json:"OLACEnabled"`
}

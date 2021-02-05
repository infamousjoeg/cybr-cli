package requests

// UpdateSafe contains the body of the Update Safe function's request
type UpdateSafe struct {
	SafeName    string `json:"SafeName,omitempty"`
	Description string `json:"Description,omitempty"`
	OLACEnabled bool   `json:"OLACEnabled,omitempty"`
	ManagingCPM string `json:"ManagingCPM,omitempty"`
}

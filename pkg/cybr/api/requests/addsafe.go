package requests

// AddSafe contains the body of the Add Safe function's request
type AddSafe struct {
	SafeName              string `json:"SafeName"`
	Description           string `json:"Description"`
	OLACEnabled           bool   `json:"OLACEnabled,omitempty"`
	ManagingCPM           string `json:"ManagingCPM"`
	NumberOfDaysRetention int    `json:"NumberOfDaysRetention"`
	AutoPurgeEnabled      bool   `json:"AutoPurgeEnabled,omitempty"`
	SafeLocation          string `json:"Location,omitempty"`
}

package responses

// GetPlatform response from getting specific account details
type GetPlatform struct {
	PlatformID string            `json:"PlatformID"`
	Details    map[string]string `json:"Details"`
	Active     bool              `json:"Active"`
}

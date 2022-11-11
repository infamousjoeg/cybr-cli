package queries

// ListPlatforms represents valid query parameters when listing platforms
type ListPlatforms struct {
	Active       bool   `query_key:"active"`
	PlatformType string `query_key:"platformType"`
	PlatformName string `query_key:"search"`
}

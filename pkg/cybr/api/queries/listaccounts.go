package queries

// ListAccounts represents valid query parameters when listing accounts
type ListAccounts struct {
	Search     string `query_key:"search"`
	SearchType string `query_key:"searchType"`
	Sort       string `query_key:"sort"`
	Offset     int    `query_key:"offset"`
	Limit      int    `query_key:"limit"`
	Filter     string `query_key:"filter"`
}

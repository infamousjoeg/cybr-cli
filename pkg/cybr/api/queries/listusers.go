package queries

// ListUsers represents valid query parameters when listing users
type ListUsers struct {
	Search string `query_key:"search"`
	Filter string `query_key:"filter"`
}

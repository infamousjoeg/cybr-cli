package queries

// ListSafeMembers represents valid query parameters when listing safe members
type ListSafeMembers struct {
	Search string `query_key:"search"`
	Sort   string `query_key:"sort"`
	Offset int    `query_key:"offset"`
	Limit  int    `query_key:"limit"`
	Filter string `query_key:"filter"`
}

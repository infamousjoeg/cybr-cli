package responses

// ListAccount response from listing accounts
type ListAccount struct {
	Value []GetAccount `json:"value"`
	Count int          `json:"count"`
}

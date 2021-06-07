package responses

// ListSafes contains an array of all safes the current user can read
type ListSafes struct {
	Safes []ListSafe `json:"value"`
}

// ListSafe contains the safe details of every safe the current user can read
// for ListSafesResponse struct
type ListSafe struct {
	SafeURLId   string `json:"SafeUrlId"`
	SafeName    string `json:"SafeName"`
	Description string `json:"Description,omitempty"`
	Location    string `json:"Location"`
}

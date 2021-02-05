package responses

// ListUsers response when listing users
type ListUsers struct {
	Users []UserResponse `json:"Users"`
	Total int            `json:"Total"`
}

// UserResponse represents one user in ListUsersResponse
type UserResponse struct {
	ID                 int                     `json:"id"`
	Username           string                  `json:"username"`
	Source             string                  `json:"source"`
	UserType           string                  `json:"userType"`
	ComponentUser      bool                    `json:"componentUser"`
	VaultAuthorization []string                `json:"vaultAuthorization"`
	Location           string                  `json:"location"`
	PersonalDetails    PersonalDetailsResponse `json:"personalDetails"`
}

// PersonalDetailsResponse represents one users personal details
type PersonalDetailsResponse struct {
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName"`
	LastName   string `json:"lastName"`
}

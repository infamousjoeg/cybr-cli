package requests

// AddSafeMember used in AddSafeMemberRequest
type AddSafeMember struct {
	MemberName               string            `json:"MemberName"`
	SearchIn                 string            `json:"SearchIn"`
	MembershipExpirationDate string            `json:"MembershipExpirationDate,omitempty"`
	Permissions              map[string]string `json:"Permissions,omitempty"`
}

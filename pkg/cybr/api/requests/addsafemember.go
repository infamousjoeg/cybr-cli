package requests

// AddSafeMember request sent for adding a member to safe with specific permissions
type AddSafeMember struct {
	Member AddSafeMemberInternal `json:"member"`
}

// AddSafeMemberInternal used in AddSafeMemberRequest
type AddSafeMemberInternal struct {
	MemberName               string               `json:"MemberName"`
	SearchIn                 string               `json:"SearchIn"`
	MembershipExpirationDate string               `json:"MembershipExpirationDate,omitempty"`
	Permissions              []PermissionKeyValue `json:"Permissions,omitempty"`
}

// PermissionKeyValue used in AddSafeMember struct
type PermissionKeyValue struct {
	Key   string `json:"Key"`
	Value bool   `json:"Value"`
}

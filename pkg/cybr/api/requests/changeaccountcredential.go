package requests

// ChangeAccountCredential only used when account is part of a group
type ChangeAccountCredential struct {
	ChangeEntireGroup bool `json:"ChangeEntireGroup,omitempty"`
}

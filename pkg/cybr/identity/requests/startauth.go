package requests

// StartAuthentication is the request to start authentication
type StartAuthentication struct {
	TenantID string `json:"TenantId,omitempty"`
	User     string `json:"User"`
	Version  string `json:"Version"`
}

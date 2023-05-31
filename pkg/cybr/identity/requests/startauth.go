package requests

// StartAuthentication is the request to start authentication
type StartAuthentication struct {
	TenantID string `json:"TenantId,omitempty"`
	Username string `json:"Username"`
	Version  string `json:"Version"`
}

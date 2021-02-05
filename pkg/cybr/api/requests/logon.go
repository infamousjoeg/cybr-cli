package requests

// Logon contains the body of the Logon function's request
type Logon struct {
	Username          string `json:"username"`
	Password          string `json:"password"`
	ConcurrentSession bool   `json:"concurrentSession,omitempty"`
}

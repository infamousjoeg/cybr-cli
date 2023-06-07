package shared

// Result contains the result response
type Result struct {
	ClientHints        ClientHint  `json:"ClientHints"`
	Version            string      `json:"Version"`
	SessionID          string      `json:"SessionId"`
	EventDescription   *string     `json:"EventDescription"`
	RetryWaitingTime   int         `json:"RetryWaitingTime"`
	SecurityImageName  *string     `json:"SecurityImageName"`
	AllowLoginMfaCache bool        `json:"AllowLoginMfaCache"`
	Challenges         []Challenge `json:"Challenges"`
	Summary            string      `json:"Summary"`
	TenantID           string      `json:"TenantId"`
}

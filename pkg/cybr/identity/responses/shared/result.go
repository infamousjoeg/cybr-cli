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
	AuthLevel          string      `json:"AuthLevel"`
	DisplayName        string      `json:"DisplayName"`
	Token              string      `json:"Token"`
	Auth               string      `json:"Auth"`
	UserID             string      `json:"UserId"`
	EmailAddress       string      `json:"EmailAddress"`
	UserDirectory      string      `json:"UserDirectory"`
	PodFqdn            string      `json:"PodFqdn"`
	User               string      `json:"User"`
	CustomerID         string      `json:"CustomerId"`
	SystemID           string      `json:"SystemId"`
	SourceDsType       string      `json:"SourceDsType"`
}

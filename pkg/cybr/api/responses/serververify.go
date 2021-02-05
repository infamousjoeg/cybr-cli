package responses

// ServerVerify contains all response data from /verify endpoint
type ServerVerify struct {
	ApplicationName       string                  `json:"ApplicationName"`
	AuthenticationMethods []AuthenticationMethods `json:"AuthenticationMethods"`
	ServerID              string                  `json:"ServerId"`
	ServerName            string                  `json:"ServerName"`
}

// AuthenticationMethods feeds into ServerVerify
type AuthenticationMethods struct {
	Enabled bool   `json:"Enabled"`
	ID      string `json:"Id"`
}

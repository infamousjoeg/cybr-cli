package shared

// ClientHint contains the client hint response
type ClientHint struct {
	PersistDefault                bool `json:"PersistDefault"`
	AllowPersist                  bool `json:"AllowPersist"`
	AllowForgotPassword           bool `json:"AllowForgotPassword"`
	EndpointAuthenticationEnabled bool `json:"EndpointAuthenticationEnabled"`
}

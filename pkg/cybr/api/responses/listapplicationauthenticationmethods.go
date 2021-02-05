package responses

// ListApplicationAuthenticationMethods contains all auth methods for a specific App ID
type ListApplicationAuthenticationMethods struct {
	Authentication []ListAuthentication `json:"authentication"`
}

// ListAuthentication contains details of the authentication method listed
type ListAuthentication struct {
	AllowInternalScripts bool   `json:"AllowInternalScripts,omitempty"`
	AppID                string `json:"AppID"`
	AuthType             string `json:"AuthType"`
	AuthValue            string `json:"AuthValue"`
	Comment              string `json:"Comment,omitempty"`
	IsFolder             bool   `json:"IsFolder,omitempty"`
	AuthID               string `json:"authID"`
}

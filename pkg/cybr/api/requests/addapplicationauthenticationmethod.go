package requests

// AddApplicationAuthentication request to add authentication method to application
type AddApplicationAuthentication struct {
	Authentication ApplicationAuthenticationMethod `json:"authentication"`
}

// ApplicationAuthenticationMethod represents an application authentication method
type ApplicationAuthenticationMethod struct {
	AuthType             string `json:"AuthType"`
	AuthValue            string `json:"AuthValue"`
	IsFolder             bool   `json:"IsFolder,omitempty"`
	AllowInternalScripts bool   `json:"AllowInternalScripts,omitempty"`
}

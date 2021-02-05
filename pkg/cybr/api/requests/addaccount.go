package requests

import "github.com/infamousjoeg/cybr-cli/pkg/cybr/api/shared"

// AddAccount request used to create an account
type AddAccount struct {
	Name                      string                  `json:"name,omitempty"`
	Address                   string                  `json:"address"`
	UserName                  string                  `json:"userName"`
	PlatformID                string                  `json:"platformId"`
	SafeName                  string                  `json:"safeName"`
	SecretType                string                  `json:"secretType"`
	Secret                    string                  `json:"secret"`
	PlatformAccountProperties map[string]string       `json:"platformAccountProperties,omitempty"`
	SecretManagement          shared.SecretManagement `json:"secretManagement,omitempty"`
}

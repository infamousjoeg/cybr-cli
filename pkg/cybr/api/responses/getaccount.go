package responses

import "github.com/infamousjoeg/cybr-cli/pkg/cybr/api/shared"

// GetAccount response from getting specific account details
type GetAccount struct {
	CategoryModificationTime  int                     `json:"categoryModificationTime"`
	ID                        string                  `json:"id"`
	Name                      string                  `json:"name"`
	Address                   string                  `json:"address"`
	UserName                  string                  `json:"userName"`
	PlatformID                string                  `json:"platformId"`
	SafeName                  string                  `json:"safeName"`
	SecretType                string                  `json:"secretType"`
	PlatformAccountProperties map[string]string       `json:"platformAccountProperties"`
	SecretManagement          shared.SecretManagement `json:"secretManagement"`
	CreatedTime               int                     `json:"createdTime"`
}

package responses

import "github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"

// AddUser ...
type AddUser struct {
	EnableUser              bool                            `json:"enableUser"`
	ChangePassOnNextLogon   bool                            `json:"changePassOnNextLogon"`
	ExpiryDate              int                             `json:"expiryDate"`
	Suspended               bool                            `json:"suspended"`
	LastSuccessfulLoginDate int                             `json:"lastSuccessfulLoginDate"`
	UnAuthorizedInterfaces  []string                        `json:"unAuthorizedInterfaces"`
	AuthenticationMethod    []string                        `json:"authenticationMethod"`
	PasswordNeverExpires    bool                            `json:"passwordNeverExpires"`
	DistinguishedName       string                          `json:"distinguishedName"`
	Description             string                          `json:"description"`
	BusinessAddress         requests.AddUserBusinessAddress `json:"businessAddress"`
	Internet                requests.AddUserInternet        `json:"internet"`
	Phones                  requests.AddUserPhones          `json:"phones"`
	PersonalDetails         requests.AddUserPersonalDetails `json:"personalDetails"`
	GroupsMembership        []interface{}                   `json:"groupsMembership"`
	ID                      int                             `json:"id"`
	Username                string                          `json:"username"`
	Source                  string                          `json:"source"`
	UserType                string                          `json:"userType"`
	ComponentUser           bool                            `json:"componentUser"`
	VaultAuthorization      []string                        `json:"vaultAuthorization"`
	Location                string                          `json:"location"`
}

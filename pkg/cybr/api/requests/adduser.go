package requests

// AddUser request that represents PAS User
type AddUser struct {
	Username               string            `json:"username"`
	UserType               string            `json:"userType"`
	InitialPassword        string            `json:"initialPassword"`
	AuthenticationMethod   []string          `json:"authenticationMethod"`
	Location               string            `json:"location"`
	UnAuthorizedInterfaces []string          `json:"unAuthorizedInterfaces"`
	ExpiryDate             int               `json:"expiryDate,omitempty"`
	VaultAuthorization     []string          `json:"vaultAuthorization"`
	EnableUser             bool              `json:"enableUser,omitempty"`
	ChangePassOnNextLogon  bool              `json:"changePassOnNextLogon,omitempty"`
	PasswordNeverExpires   bool              `json:"passwordNeverExpires"`
	DistinguishedName      string            `json:"distinguishedName"`
	Description            string            `json:"description"`
	BusinessAddress        map[string]string `json:"businessAddress,omitempty"`
	Internet               map[string]string `json:"internet,omitempty"`
	Phones                 map[string]string `json:"phones,omitempty"`
	PersonalDetails        map[string]string `json:"personalDetails,omitempty"`
}

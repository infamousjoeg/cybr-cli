package requests

// AddUser the name of the struct should be the same name as the method it is implemented in
type AddUser struct {
	Username               string                 `json:"username"`
	UserType               string                 `json:"userType"`
	InitialPassword        string                 `json:"initialPassword"`
	AuthenticationMethod   []string               `json:"authenticationMethod"`
	Location               string                 `json:"location"`
	UnAuthorizedInterfaces []string               `json:"unAuthorizedInterfaces"`
	ExpiryDate             int                    `json:"expiryDate,omitempty"`
	VaultAuthorization     []string               `json:"vaultAuthorization"`
	EnableUser             bool                   `json:"enableUser"`
	ChangePassOnNextLogon  bool                   `json:"changePassOnNextLogon"`
	PasswordNeverExpires   bool                   `json:"passwordNeverExpires"`
	DistinguishedName      string                 `json:"distinguishedName"`
	Description            string                 `json:"description"`
	BusinessAddress        AddUserBusinessAddress `json:"businessAddress"`
	Internet               AddUserInternet        `json:"internet"`
	Phones                 AddUserPhones          `json:"phones"`
	PersonalDetails        AddUserPersonalDetails `json:"personalDetails"`
}

// AddUserBusinessAddress ...
type AddUserBusinessAddress struct {
	WorkStreet  string `json:"workStreet"`
	WorkCity    string `json:"workCity"`
	WorkState   string `json:"workState"`
	WorkZip     string `json:"workZip"`
	WorkCountry string `json:"workCountry"`
}

// AddUserInternet ...
type AddUserInternet struct {
	HomePage      string `json:"homePage"`
	HomeEmail     string `json:"homeEmail"`
	BusinessEmail string `json:"businessEmail"`
	OtherEmail    string `json:"otherEmail"`
}

// AddUserPhones ...
type AddUserPhones struct {
	HomeNumber     string `json:"homeNumber"`
	BusinessNumber string `json:"businessNumber"`
	CellularNumber string `json:"cellularNumber"`
	FaxNumber      string `json:"faxNumber"`
	PagerNumber    string `json:"pagerNumber"`
}

// AddUserPersonalDetails ...
type AddUserPersonalDetails struct {
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	Zip          string `json:"zip"`
	Country      string `json:"country"`
	Title        string `json:"title"`
	Organization string `json:"organization"`
	Department   string `json:"department"`
	Profession   string `json:"profession"`
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
}

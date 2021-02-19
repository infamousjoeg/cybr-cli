package responses

// AddUser response
type AddUser struct {
	EnableUser              bool            `json:"enableUser"`
	ChangePassOnNextLogon   bool            `json:"changePassOnNextLogon"`
	ExpiryDate              int             `json:"expiryDate"`
	Suspended               bool            `json:"suspended"`
	LastSuccessfulLoginDate int             `json:"lastSuccessfulLoginDate"`
	UnAuthorizedInterfaces  []string        `json:"unAuthorizedInterfaces"`
	AuthenticationMethod    []string        `json:"authenticationMethod"`
	PasswordNeverExpires    bool            `json:"passwordNeverExpires"`
	DistinguishedName       string          `json:"distinguishedName"`
	Description             string          `json:"description"`
	BusinessAddress         BusinessAddress `json:"businessAddress,omitempty"`
	Internet                Internet        `json:"internet,omitempty"`
	Phones                  Phones          `json:"phones,omitempty"`
	PersonalDetails         PersonalDetails `json:"personalDetails,omitempty"`
	ID                      int             `json:"id"`
	Username                string          `json:"username"`
	Source                  string          `json:"source"`
	UserType                string          `json:"userType"`
	ComponentUser           bool            `json:"componentUser"`
	VaultAuthorization      []string        `json:"vaultAuthorization"`
	Location                string          `json:"location"`
}

// BusinessAddress of user
type BusinessAddress struct {
	WorkStreet  string `json:"workStreet"`
	WorkCity    string `json:"workCity"`
	WorkState   string `json:"workState"`
	WorkZip     string `json:"workZip"`
	WorkCountry string `json:"workCountry"`
}

// Internet of user
type Internet struct {
	HomePage      string `json:"homePage"`
	HomeEmail     string `json:"homeEmail"`
	BusinessEmail string `json:"businessEmail"`
	OtherEmail    string `json:"otherEmail"`
}

// Phones of user
type Phones struct {
	HomeNumber     string `json:"homeNumber"`
	BusinessNumber string `json:"businessNumber"`
	CellularNumber string `json:"cellularNumber"`
	FaxNumber      string `json:"faxNumber"`
	PagerNumber    string `json:"pagerNumber"`
}

// PersonalDetails of user
type PersonalDetails struct {
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

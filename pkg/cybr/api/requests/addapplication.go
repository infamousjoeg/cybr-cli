package requests

// AddApplication request for adding application
type AddApplication struct {
	Application Application `json:"application"`
}

// Application is used when adding an application
type Application struct {
	AppID               string `json:"AppID"`
	Description         string `json:"Description,omitempty"`
	Location            string `json:"Location,omitempty"`
	AccessPermittedFrom int    `json:"AccessPermittedFrom"`
	AccessPermittedTo   int    `json:"AccessPermittedTo"`
	ExpirationDate      string `json:"ExpirationDate,omitempty"`
	Disabled            string `json:"Disabled,omitempty"`
	BusinessOwnerFName  string `json:"BusinessOwnerFName,omitempty"`
	BusinessOwnerLName  string `json:"BusinessOwnerLName,omitempty"`
	BusinessOwnerEmail  string `json:"BusinessOwnerEmail,omitempty"`
	BusinessOwnerPhone  string `json:"BusinessOwnerPhone,omitempty"`
}

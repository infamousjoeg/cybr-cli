package responses

import "time"

//ListApplications contains all applications and app data returned
type ListApplications struct {
	Application []ListApplication `json:"application"`
}

// ListApplication contains all specific application data for ListApplicationsResponse
type ListApplication struct {
	AccessPermittedFrom                     int       `json:"AccessPermittedFrom"`
	AccessPermittedTo                       int       `json:"AccessPermittedTo"`
	AllowExtendedAuthenticationRestrictions bool      `json:"AllowExtendedAuthenticationRestrict"`
	AppID                                   string    `json:"AppID"`
	BusinessOwnerEmail                      string    `json:"BusinessOwnerEmail"`
	BusinessOwnerFName                      string    `json:"BusinessOwnerFName"`
	BusinessOwnerLName                      string    `json:"BusinessOwnerLName"`
	BusinessOwnerPhone                      string    `json:"BusinessOwnerPhone"`
	Description                             string    `json:"Description"`
	Disabled                                bool      `json:"Disabled"`
	ExpirationDate                          time.Time `json:"ExpirationDate"`
	Location                                string    `json:"Location"`
}

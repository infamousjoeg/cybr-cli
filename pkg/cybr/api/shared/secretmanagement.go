package shared

// SecretManagement used in getting and setting accounts
type SecretManagement struct {
	AutomaticManagementEnabled bool   `json:"automaticManagementEnabled"`
	Status                     string `json:"status"`
	ManualManagementReason     string `json:"manualManagementReason,omitempty"`
	LastModifiedTime           int    `json:"lastModifiedTime,omitempty"`
}

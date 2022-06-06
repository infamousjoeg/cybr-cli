package responses

// ListPlatforms contains an array of all platforms
type ListPlatforms struct {
	Platforms []ListPlatform `json:"Platforms"`
}

// ListPlatform contains the platform details of every platform
// for ListPlatformsResponse struct
type ListPlatform struct {
	General                   General                   `json:"general"`
	Properties                Properties                `json:"properties"`
	LinkedAccounts            []LinkedAccounts          `json:"linkedAccounts,omitempty"`
	CredentialsManagement     []CredentialsManagement   `json:"creditentialsManagement"`
	SessionManagement         SessionManagement         `json:"sessionManagement"`
	PrivilegedAccessWorkflows PrivilegedAccessWorkflows `json:"privilegedAccessWorkflows"`
}

// General contains the general details of a platform
type General struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	SystemType     string `json:"systemType"`
	Active         bool   `json:"active"`
	Description    string `json:"description,omitempty"`
	PlatformBaseID string `json:"platformBaseId"`
	PlatformType   string `json:"platformType"`
}

// Properties contains the properties of a platform
type Properties struct {
	Required []RequiredProperties `json:"required,omitempty"`
	Optional []OptionalProperties `json:"optional,omitempty"`
}

// RequiredProperties contains the required properties of a platform
type RequiredProperties struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// OptionalProperties contains the optional properties of a platform
type OptionalProperties struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// LinkedAccounts contains the linked accounts of a platform
type LinkedAccounts struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName,omitempty"`
}

// CredentialsManagement contains the credentials management of a platform
type CredentialsManagement struct {
	AllowedSafes                          string `json:"allowedSafes"`
	AllowManualChange                     bool   `json:"allowManualChange"`
	PerformPeriodicChange                 bool   `json:"performPeriodicChange"`
	RequirePasswordChangeEveryXDays       int    `json:"requirePasswordChangeEveryXDays"`
	AllowManualVerification               bool   `json:"allowManualVerification"`
	PerformPeriodicVerification           bool   `json:"performPeriodicVerification"`
	RequirePasswordVerificationEveryXDays int    `json:"requirePasswordVerificationEveryXDays"`
	AllowManualReconciliation             bool   `json:"allowManualReconciliation"`
	AutomaticReconcileWhenUnsynched       bool   `json:"automaticReconcileWhenUnsynched"`
}

// SessionManagement contains the session management of a platform
type SessionManagement struct {
	RequirePrivilegedSessionMonitoringAndIsolation bool   `json:"requirePrivilegedSessionMonitoringAndIsolation"`
	RecordAndSaveSessionActivity                   bool   `json:"recordAndSaveSessionActivity"`
	PSMServerID                                    string `json:"psmServerId,omitempty"`
}

// PrivilegedAccessWorkflows contains the privileged access workflows of a platform
type PrivilegedAccessWorkflows struct {
	RequireDualControlPasswordAccessApproval bool `json:"requireDualControlPasswordAccessApproval"`
	EnforceCheckinCheckoutExclusiveAccess    bool `json:"enforceCheckinCheckoutExclusiveAccess"`
	EnforceOnetimePasswordAccess             bool `json:"enforceOnetimePasswordAccess"`
}

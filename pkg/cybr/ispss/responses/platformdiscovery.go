package responses

// PlatformDiscovery is the response from the ISPSS API for the platform discovery endpoint
type PlatformDiscovery struct {
	SecretsManager         Data `json:"secrets_manager,omitempty"`
	Pcloud                 Data `json:"pcloud,omitempty"`
	SecretsHub             Data `json:"secrets_hub,omitempty"`
	IdaptiveRiskAnalytics  Data `json:"idaptive_risk_analytics,omitempty"`
	ComponentManager       Data `json:"component_manager,omitempty"`
	IdentityUserPortal     Data `json:"identity_user_portal,omitempty"`
	Cem                    Data `json:"cem,omitempty"`
	IdentityAdministration Data `json:"identity_administration,omitempty"`
	CloudOnboarding        Data `json:"cloud_onboarding,omitempty"`
	Analytics              Data `json:"analytics,omitempty"`
	Sca                    Data `json:"sca,omitempty"`
	Jit                    Data `json:"jit,omitempty"`
	SessionMonitoring      Data `json:"session_monitoring,omitempty"`
	Audit                  Data `json:"audit,omitempty"`
}

// Data is the data for each platform
type Data struct {
	UI        string `json:"ui,omitempty"`
	API       string `json:"api,omitempty"`
	Bootstrap string `json:"bootstrap,omitempty"`
	Region    string `json:"region,omitempty"`
}

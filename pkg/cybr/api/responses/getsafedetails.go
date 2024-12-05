package responses

type GetSafeDetails struct {
	Accounts                  []AccountInfo `json:"accounts,omitempty"`
	AutoPurgeEnabled          bool          `json:"autoPurgeEnabled,omitempty"`
	CreationTime              int64         `json:"creationTime,omitempty"`
	Creator                   AccountInfo   `json:"creator,omitempty"`
	Description               string        `json:"description,omitempty"`
	IsExpiredMember           bool          `json:"isExpiredMember,omitempty"`
	LastModificationTime      int64         `json:"lastModificationTime,omitempty"`
	Location                  string        `json:"location,omitempty"`
	ManagingCPM               string        `json:"managingCPM,omitempty"`
	NumberOfDaysRetention     int           `json:"numberOfDaysRetention,omitempty"`
	NumberOfVersionsRetention int           `json:"numberOfVersionsRetention,omitempty"`
	OlacEnabled               bool          `json:"olacEnabled,omitempty"`
	SafeName                  string        `json:"safeName,omitempty"`
	SafeNumber                int           `json:"safeNumber,omitempty"`
	SafeURLID                 string        `json:"safeUrlId,omitempty"`
}

type AccountInfo struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

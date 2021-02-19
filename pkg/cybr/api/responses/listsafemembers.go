package responses

// ListSafeMembers contains data of all members of a specific safe
type ListSafeMembers struct {
	Members []Members `json:"members"`
}

// Members contains all safe member username/group name and their permissions
type Members struct {
	Permissions Permissions `json:"Permissions"`
	Username    string      `json:"UserName"`
}

// Permissions contains the permissions of each safe member
type Permissions struct {
	Add                 bool `json:"Add"`
	AddRenameFolder     bool `json:"AddRenameFolder"`
	BackupSafe          bool `json:"BackupSafe"`
	Delete              bool `json:"Delete"`
	DeleteFolder        bool `json:"DeleteFolder"`
	ListContent         bool `json:"ListContent"`
	ManageSafe          bool `json:"ManageSafe"`
	ManageSafeMembers   bool `json:"ManageSafeMembers"`
	MoveFilesAndFolders bool `json:"MoveFilesAndFolders"`
	Rename              bool `json:"Rename"`
	RestrictedRetrieve  bool `json:"RestrictedRetrieve"`
	Retrieve            bool `json:"Retrieve"`
	Unlock              bool `json:"Unlock"`
	Update              bool `json:"Update"`
	UpdateMetadata      bool `json:"UpdateMetadata"`
	ValidateSafeContent bool `json:"ValidateSafeContent"`
	ViewAudit           bool `json:"ViewAudit"`
	ViewMembers         bool `json:"ViewMembers"`
}

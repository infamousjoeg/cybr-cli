package responses

// ListSafeMembers contains data of all members of a specific safe
type ListSafeMembers struct {
	Members []Members `json:"value"`
	Count   int       `json:"count"`
}

// Members contains all safe member username/group name and their permissions
type Members struct {
	SafeName                  string      `json:"safeName"`
	SafeNumber                int         `json:"safeNumber"`
	MemberID                  int         `json:"memberId"`
	MemberName                string      `json:"memberName"`
	MemberType                string      `json:"memberType"`
	IsExpiredMembershipEnable bool        `json:"isExpiredMembershipEnable"`
	IsPredefinedUser          bool        `json:"isPredefinedUser"`
	Permissions               Permissions `json:"Permissions"`
}

// Permissions contains the permissions of each safe member
type Permissions struct {
	UseAccounts                            bool `json:"UseAccounts"`
	RetrieveAccounts                       bool `json:"RetrieveAccounts"`
	ListAccounts                           bool `json:"ListAccounts"`
	AddAccounts                            bool `json:"AddAccounts"`
	UpdateAccountContent                   bool `json:"UpdateAccountContent"`
	UpdateAccountProperties                bool `json:"UpdateAccountProperties"`
	InitiateCPMAccountManagementOperations bool `json:"InitiateCPMAccountManagementOperations"`
	SpecifyNextAccountContent              bool `json:"SpecifyNextAccountContent"`
	RenameAccounts                         bool `json:"RenameAccounts"`
	DeleteAccounts                         bool `json:"DeleteAccounts"`
	UnlockAccounts                         bool `json:"UnlockAccounts"`
	ManageSafe                             bool `json:"ManageSafe"`
	ManageSafeMembers                      bool `json:"ManageSafeMembers"`
	BackupSafe                             bool `json:"BackupSafe"`
	ViewAuditLog                           bool `json:"ViewAuditLog"`
	ViewSafeMembers                        bool `json:"ViewSafeMembers"`
	AccessWithoutConfirmation              bool `json:"AccessWithoutConfirmation"`
	CreateFolders                          bool `json:"CreateFolders"`
	DeleteFolders                          bool `json:"DeleteFolders"`
	MoveAccountsAndFolders                 bool `json:"MoveAccountsAndFolders"`
	RequestsAuthorizationLevel1            bool `json:"RequestsAuthorizationLevel1"`
	RequestsAuthorizationLevel2            bool `json:"RequestsAuthorizationLevel2"`
}

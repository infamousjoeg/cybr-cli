package requests

type UpdateSafeMember struct {
	MembershipExpirationDate int64       `json:"membershipExpirationDate,omitempty"`
	Permissions              Permissions `json:"permissions,omitempty"`
}
type Permissions struct {
	UseAccounts                            bool `json:"useAccounts,omitempty"`
	RetrieveAccounts                       bool `json:"retrieveAccounts,omitempty"`
	ListAccounts                           bool `json:"listAccounts,omitempty"`
	AddAccounts                            bool `json:"addAccounts,omitempty"`
	UpdateAccountContent                   bool `json:"updateAccountContent,omitempty"`
	UpdateAccountProperties                bool `json:"updateAccountProperties,omitempty"`
	InitiateCPMAccountManagementOperations bool `json:"initiateCPMAccountManagementOperations,omitempty"`
	SpecifyNextAccountContent              bool `json:"specifyNextAccountContent,omitempty"`
	RenameAccounts                         bool `json:"renameAccounts,omitempty"`
	DeleteAccounts                         bool `json:"deleteAccounts,omitempty"`
	UnlockAccounts                         bool `json:"unlockAccounts,omitempty"`
	ManageSafe                             bool `json:"manageSafe,omitempty"`
	ManageSafeMembers                      bool `json:"manageSafeMembers,omitempty"`
	BackupSafe                             bool `json:"backupSafe,omitempty"`
	ViewAuditLog                           bool `json:"viewAuditLog,omitempty"`
	ViewSafeMembers                        bool `json:"viewSafeMembers,omitempty"`
	AccessWithoutConfirmation              bool `json:"accessWithoutConfirmation,omitempty"`
	CreateFolders                          bool `json:"createFolders,omitempty"`
	DeleteFolders                          bool `json:"deleteFolders,omitempty"`
	MoveAccountsAndFolders                 bool `json:"moveAccountsAndFolders,omitempty"`
	RequestsAuthorizationLevel1            bool `json:"requestsAuthorizationLevel1,omitempty"`
	RequestsAuthorizationLevel2            bool `json:"requestsAuthorizationLevel2,omitempty"`
}

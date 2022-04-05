package api

import (
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

// GetRolePermissions assigns pre-defined safe permissions for new safe member
func GetRolePermissions(role string) ([]requests.PermissionKeyValue, error) {
	var permissions []requests.PermissionKeyValue

	// Set permissions variable to pre-define safe permissions based on role given
	if role == "ApplicationIdentity" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "RetrieveAccounts",
				Value: true,
			},
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "AIMWebService" {
		permissions = []requests.PermissionKeyValue{}

		return permissions, nil
	}

	if role == "AppProvider" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "RetrieveAccounts",
				Value: true,
			},
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "BreakGlass" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "UseAccounts",
				Value: true,
			},
			{
				Key:   "RetrieveAccounts",
				Value: true,
			},
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AddAccounts",
				Value: true,
			},
			{
				Key:   "UpdateAccountContent",
				Value: true,
			},
			{
				Key:   "UpdateAccountProperties",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "SpecifyNextAccountContent",
				Value: true,
			},
			{
				Key:   "RenameAccounts",
				Value: true,
			},
			{
				Key:   "DeleteAccounts",
				Value: true,
			},
			{
				Key:   "UnlockAccounts",
				Value: true,
			},
			{
				Key:   "ManageSafe",
				Value: true,
			},
			{
				Key:   "ManageSafeMembers",
				Value: true,
			},
			{
				Key:   "BackupSafe",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
			{
				Key:   "CreateFolders",
				Value: true,
			},
			{
				Key:   "DeleteFolders",
				Value: true,
			},
			{
				Key:   "MoveAccountsAndFolders",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "AccountProvisioner" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AddAccounts",
				Value: true,
			},
			{
				Key:   "UpdateAccountContent",
				Value: true,
			},
			{
				Key:   "UpdateAccountProperties",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "RenameAccounts",
				Value: true,
			},
			{
				Key:   "DeleteAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "CreateFolders",
				Value: true,
			},
			{
				Key:   "DeleteFolders",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "CPDeployer" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AddAccounts",
				Value: true,
			},
			{
				Key:   "UpdateAccountProperties",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "PasswordScheduler" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "UnlockAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "SafeManager" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ManageSafe",
				Value: true,
			},
			{
				Key:   "ManageSafeMembers",
				Value: true,
			},
			{
				Key:   "BackupSafe",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "ComponentOrchestrator" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AddAccounts",
				Value: true,
			},
			{
				Key:   "UpdateAccountProperties",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "AccessWithoutConfirmation",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "EndUser" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "UseAccounts",
				Value: true,
			},
			{
				Key:   "RetrieveAccounts",
				Value: true,
			},
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "VaultAdmin" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "ApproverLevel1" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "RequestsAuthorizationLevel1",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "ApproverLevel2" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "RequestsAuthorizationLevel2",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "AddAccounts",
				Value: true,
			},
			{
				Key:   "UpdateAccountContent",
				Value: true,
			},
			{
				Key:   "UpdateAccountProperties",
				Value: true,
			},
			{
				Key:   "InitiateCPMAccountManagementOperations",
				Value: true,
			},
			{
				Key:   "RenameAccounts",
				Value: true,
			},
			{
				Key:   "DeleteAccounts",
				Value: true,
			},
			{
				Key:   "UnlockAccounts",
				Value: true,
			},
			{
				Key:   "ManageSafe",
				Value: true,
			},
			{
				Key:   "ManageSafeMembers",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
			{
				Key:   "CreateFolders",
				Value: true,
			},
			{
				Key:   "DeleteFolders",
				Value: true,
			},
			{
				Key:   "MoveAccountsAndFolders",
				Value: true,
			},
		}

		return permissions, nil
	}

	if role == "Auditor" {
		permissions = []requests.PermissionKeyValue{
			{
				Key:   "ListAccounts",
				Value: true,
			},
			{
				Key:   "ViewAuditLog",
				Value: true,
			},
			{
				Key:   "ViewSafeMembers",
				Value: true,
			},
		}

		return permissions, nil
	}

	return permissions, fmt.Errorf("Unknown role value")
}

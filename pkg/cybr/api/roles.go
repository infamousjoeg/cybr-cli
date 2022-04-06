package api

import (
	"fmt"

	"github.com/infamousjoeg/cybr-cli/pkg/cybr/api/requests"
)

// GetRolePermissions assigns pre-defined safe permissions for new safe member
func GetRolePermissions(role string) ([]requests.PermissionKeyValue, error) {
	var permissions []requests.PermissionKeyValue

	// Set permissions variable to pre-define safe permissions based on role given
	if role == "BreakGlass" ||
		role == "EndUser" {
		permissions = []requests.PermissionKeyValue{
			{Key: "UseAccounts", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "ApplicationIdentity" ||
		role == "AppProvider" ||
		role == "EndUser" {
		permissions = []requests.PermissionKeyValue{
			{Key: "RetrieveAccounts", Value: true},
		}
	} else if role != "AIMWebService" &&
		role != "SafeManager" {
		permissions = []requests.PermissionKeyValue{
			{Key: "ListAccounts", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "AccountProvisioner" ||
		role == "CPDeployer" ||
		role == "ComponentOrchestrator" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "AddAccounts", Value: true},
			{Key: "UpdateAccountProperties", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "AccountProvisioner" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "UpdateAccountContent", Value: true},
			{Key: "RenameAccounts", Value: true},
			{Key: "DeleteAccounts", Value: true},
			{Key: "CreateFolders", Value: true},
			{Key: "DeleteFolders", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "AccountProvisioner" ||
		role == "CPDeployer" ||
		role == "PasswordScheduler" ||
		role == "ComponentOrchestrator" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "InitiateCPMAccountManagementOperations", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "PasswordScheduler" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "UnlockAccounts", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "SafeManager" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "ManageSafe", Value: true},
			{Key: "ManageSafeMembers", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "SafeManager" {
		permissions = []requests.PermissionKeyValue{
			{Key: "BackupSafe", Value: true},
		}
	} else if role != "AIMWebService" &&
		role != "ApplicationIdentity" &&
		role != "AppProvider" {
		permissions = []requests.PermissionKeyValue{
			{Key: "ViewAuditLog", Value: true},
		}
	} else if role != "AIMWebService" &&
		role != "ApplicationIdentity" &&
		role != "ComponentOrchestrator" {
		permissions = []requests.PermissionKeyValue{
			{Key: "ViewSafeMembers", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "ApplicationIdentity" ||
		role == "CPDeployer" ||
		role == "PasswordScheduler" ||
		role == "SafeManager" ||
		role == "ComponentOrchestrator" {
		permissions = []requests.PermissionKeyValue{
			{Key: "AccessWithoutConfirmation", Value: true},
		}
	} else if role == "BreakGlass" ||
		role == "APIAutomation" {
		permissions = []requests.PermissionKeyValue{
			{Key: "MoveAccountsAndFolders", Value: true},
		}
	} else if role == "BreakGlass" {
		permissions = []requests.PermissionKeyValue{
			{Key: "SpecifyNextAccountContent", Value: true},
		}
	} else if role == "ApproverLevel1" {
		permissions = []requests.PermissionKeyValue{
			{Key: "RequestsAuthorizationLevel1", Value: true},
		}
	} else if role == "ApproverLevel2" {
		permissions = []requests.PermissionKeyValue{
			{Key: "RequestsAuthorizationLevel2", Value: true},
		}
	} else {
		return permissions, fmt.Errorf("Unknown role value")
	}

	return permissions, nil
}

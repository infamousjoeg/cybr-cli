package api

import (
	"fmt"
	"strings"
)

// The content will look like
// port=something, sp
func keyValueStringToMap(content string) (map[string]string, error) {
	if content == "" {
		return nil, nil
	}

	if !strings.Contains(content, "=") {
		return nil, fmt.Errorf("Invalid platform prop content. The provided content does not container a '='")
	}

	m := make(map[string]string)

	// TODO: Gotta be a better way to do this
	replaceWith := "^||||^"

	// If the address or property contains a `\,` then replace
	content = strings.ReplaceAll(content, "\\,", replaceWith)
	props := strings.Split(content, ",")
	for _, prop := range props {
		if !strings.Contains(prop, "=") {
			return nil, fmt.Errorf("Property '%s' is invalid because it does not contain a '=' to seperate key from value", prop)
		}
		kvs := strings.SplitN(prop, "=", 2)
		key := strings.Trim(kvs[0], " ")
		value := strings.Trim(strings.ReplaceAll(kvs[1], replaceWith, ","), " ")
		m[key] = value
	}

	return m, nil
}

// GetRolePermissions assigns pre-defined safe permissions for new safe member
func GetRolePermissions(role string) (map[string]string, error) {
	var permissions map[string]string

	// Set permissions variable to pre-define safe permissions based on role given
	if role == "BreakGlass" {
		permissions, err := keyValueStringToMap("UseAccounts=true,RetrieveAccounts=true,ListAccounts=true,AddAccounts=true,UpdateAccountContent=true,UpdateAccountProperties=true,InitiateCPMAccountManagementOperations=true,SpecifyNextAccountContent=true,RenameAccounts=true,DeleteAccounts=true,UnlockAccounts=true,ManageSafe=true,ManageSafeMembers=true,BackupSafe=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true,CreateFolders=true,DeleteFolders=true,MoveAccountsAndFolders=true,RequestsAuthorizationLevel1=true,RequestsAuthorizationLevel2=false")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "VaultAdmin" {
		permissions, err := keyValueStringToMap("ListAccounts=true,ViewAuditLog=true,ViewSafeMembers=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "SafeManager" {
		permissions, err := keyValueStringToMap("ManageSafe=true,ManageSafeMembers=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "EndUser" {
		permissions, err := keyValueStringToMap("UseAccounts=true,RetrieveAccounts=true,ListAccounts=true,ViewAuditLog=true,ViewSafeMembers=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "Auditor" {
		permissions, err := keyValueStringToMap("ListAccounts=true,ViewAuditLog=true,ViewSafeMembers=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "AIMWebService" {
		permissions, err := keyValueStringToMap("")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "AppProvider" {
		permissions, err := keyValueStringToMap("RetrieveAccounts=true,ListAccounts=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "ApplicationIdentity" {
		permissions, err := keyValueStringToMap("RetrieveAccounts=true,ListAccounts=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "AccountProvisioner" {
		permissions, err := keyValueStringToMap("ListAccounts=true,AddAccounts=true,UpdateAccountProperties=true,InitiateCPMAccountManagementOperations=true,DeleteAccounts=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "CPDeployer" {
		permissions, err := keyValueStringToMap("ListAccounts=true,AddAccounts=true,UpdateAccountProperties=true,InitiateCPMAccountManagementOperations=true,ManageSafeMembers=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "ComponentOrchestrator" {
		permissions, err := keyValueStringToMap("ListAccounts=true,AddAccounts=true,UpdateAccountProperties=true,InitiateCPMAccountManagementOperations=true,ViewAuditLog=true,AccessWithoutConfirmation=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "APIAutomation" {
		permissions, err := keyValueStringToMap("ListAccounts=true,AddAccounts=true,UpdateAccountContent=true,UpdateAccountProperties=true,InitiateCPMAccountManagementOperations=true,RenameAccounts=true,DeleteAccounts=true,UnlockAccounts=true,ManageSafe=true,ManageSafeMembers=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true,CreateFolders=true,DeleteFolders=true,MoveAccountsAndFolders=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "PasswordScheduler" {
		permissions, err := keyValueStringToMap("ListAccounts=true,InitiateCPMAccountManagementOperations=true,ViewAuditLog=true,ViewSafeMembers=true,AccessWithoutConfirmation=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "ApproverLevel1" {
		permissions, err := keyValueStringToMap("ListAccounts=true,ViewAuditLog=true,ViewSafeMembers=true,RequestsAuthorizationLevel1=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else if role == "ApproverLevel2" {
		permissions, err := keyValueStringToMap("ListAccounts=true,ViewAuditLog=true,ViewSafeMembers=true,RequestsAuthorizationLevel2=true")
		if err != nil {
			return nil, err
		}
		return permissions, nil
	} else {
		return permissions, fmt.Errorf("Unknown role value")
	}
}

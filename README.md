# cybr-cli <!-- omit in toc -->
@CyberArk Privileged Access Security (PAS) REST API Client Library

[![cybr-cli CI](https://github.com/infamousjoeg/cybr-cli/workflows/cybr-cli%20CI/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3A%22cybr-cli+CI%22) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=infamousjoeg_pas-api-go&metric=alert_status)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ALint) [![CodeQL](https://github.com/infamousjoeg/cybr-cli/workflows/CodeQL/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ACodeQL)

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
  - [Command-Line Interface (CLI)](#command-line-interface-cli)
    - [logon](#logon)
    - [logoff](#logoff)
    - [safes](#safes)
      - [list](#list)
      - [member list](#member-list)
    - [applications](#applications)
      - [list](#list-1)
      - [methods list](#methods-list)
    - [version](#version)
    - [help](#help)
  - [Install from Source](#install-from-source)
  - [Docker Container](#docker-container)
    - [Run Container Indefinitely](#run-container-indefinitely)
    - [Run Container Ephemerally (Recommended)](#run-container-ephemerally-recommended)
      - [One-Time Use](#one-time-use)
      - [One-Time Use w/ Saved Config](#one-time-use-w-saved-config)
      - [Using with jq](#using-with-jq)
  - [Application](#application)
    - [Import into project](#import-into-project)
    - [Logon to the PAS REST API Web Service](#logon-to-the-pas-rest-api-web-service)
    - [Call functions by referencing `pasapi` and "dot-referencing"](#call-functions-by-referencing-pasapi-and-dot-referencing)
- [Required Environment Variables](#required-environment-variables)
- [Testing](#testing)
  - [Successful Output](#successful-output)

## Usage

### Command-Line Interface (CLI)

#### logon

```shell
$ cybr logon -u username -a cyberark-or-ldap -b https://pvwa.example.com
```

__Required Options:__
* `-u` or `--username`
* `-a` or `--auth-type`
* `-b` or `--base-url`

Logon to the PAS REST API as the username you provide using the authentication method you choose. At this time, only `cyberark` and `ldap` authentication methods are supported.

Upon successful logon, a file will be created in your user's home directory at `.cybr/config`. It is an encoded file that cannot be read in plain-text. This holds your current session information.

#### logoff

```shell
$ cybr logoff
```

Logoff the PAS REST API as the username you provided during logon.

Upon successful logoff, the config file located in your user's home directory at `.cybr/config` will be removed and the session token stored within will be expired.

#### safes

```shell
$ cybr safes
```

List all safes the username you are logged on as has access to read.

##### list

```shell
$ cybr safes list
```

List all safes the username you are logged on as has access to read.

##### member list

```shell
$ cybr safes member list -s SafeName
```

__Required Option:__ `-s` or `--safe`

List all safe members on the safe given.

#### applications

```shell
$ cybr applications
```

List all applications the username you are logged on as has access to read.

```shell
$ cybr applications -l \\Applications
```

__Optional Option:__ `-l` or `--location`

List only applications located within \Applications the username you are logged on as has access to read.

##### list

```shell
$ cybr applications list
```

__Optional Option:__ `-l` or `--location`

List all applications the username you are logged on as has access to read.

##### methods list

```shell
$ cybr applications methods list -a AppID
```

__Required Option:__ `-a` or `--app-id`

List all authentication methods configured for the application identity given.

#### version

```shell
$ cybr version
```

Displays the current version of the `cybr` CLI.

#### help

```shell
$ cybr help [command]
```

Displays help text for the `cybr` CLI.  If an optional `[command]` is provided, the help text for that command will be displayed instead.

### Install from Source

```shell
$ git clone https://github.com/infamousjoeg/pas-api-go.git
$ ./install
$ cybr help
cybr is a command-line interface utility created by CyberArk that
wraps the PAS REST API and eases the user experience for automators
and automation to easily interact with CyberArk Privileged Access
Security.

Usage:
  cybr [command]

Available Commands:
  applications Applications actions for PAS REST API
  help         Help about any command
  logoff       Logoff the PAS REST API
  logon        Logon to PAS REST API
  safes        Safe actions for PAS REST API
  version      Display current version

Flags:
  -h, --help   help for cybr

Use "cybr [command] --help" for more information about a command.
```

### Docker Container

The `cybr` CLI is also available within a lightweight container for ephemeral use, if necessary.

#### Run Container Indefinitely

```shell
$ docker run --name cybr-cli -d --restart always \
  --entrypoint sleep \
  nfmsjoeg/cybr-cli:latest \
  infinity
```

Running indefinitely allows you to stay inside the container with the `cybr` CLI.

```shell
$ docker exec -it cybr-cli /bin/bash
cybr@6e1c196a84a6:/app$ cybr version
cybr v0.0.2-alpha
```

#### Run Container Ephemerally (Recommended)

##### One-Time Use

```shell
$ docker run --rm -it nfmsjoeg/cybr-cli:latest /bin/bash
cybr@6e1c196a84a6:/app$
```

##### One-Time Use w/ Saved Config

```shell
$ docker run --rm -it \
  -v cybr-cli:/home/cybr/.cybr \
  nfmsjoeg/cybr-cli:latest \
  logon -u username -a cyberark -b https://pvwa.example.com
Enter password:
Successfully logged onto PAS as user jgarcia.
$ docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
$ docker run --rm \
  -v cybr-cli:/home/cybr/.cybr \
  nfmsjoeg/cybr-cli:latest \
  safes list
{
    "Safes": [
        {
            "SafeUrlId": "VaultInternal",
            "SafeName": "VaultInternal",
            "Location": "\\"
        },
        {
            "SafeUrlId": "Notification%20Engine",
            "SafeName": "Notification Engine",
            "Location": "\\"
        },
        {
            "SafeUrlId": "PVWAReports",
            "SafeName": "PVWAReports",
            "Location": "\\"
        },
        {
            "SafeUrlId": "PVWATicketingSystem",
            "SafeName": "PVWATicketingSystem",
            "Location": "\\"
        },
        {
            "SafeUrlId": "PVWAPublicData",
            "SafeName": "PVWAPublicData",
            "Location": "\\"
        },
        {
            "SafeUrlId": "PasswordManager",
            "SafeName": "PasswordManager",
            "Location": "\\"
        },
        {
            "SafeUrlId": "PasswordManager_Pending",
            "SafeName": "PasswordManager_Pending",
            "Location": "\\"
        },
        {
            "SafeUrlId": "AccountsFeedADAccounts",
            "SafeName": "AccountsFeedADAccounts",
            "Location": "\\"
        },
        {
            "SafeUrlId": "AccountsFeedDiscoveryLogs",
            "SafeName": "AccountsFeedDiscoveryLogs",
            "Location": "\\"
        }
	]
}
```

##### Using with jq

You can also pipe output to `jq` [[download]](https://stedolan.github.io/jq/) to easily parse the returned JSON responses:

```shell
$ docker run --rm \
  -v cybr-cli:/home/cybr/.cybr \
  nfmsjoeg/cybr-cli:latest \
  safes list | jq '.Safes[] | {SafeName}'
{
  "SafeName": "VaultInternal"
}
{
  "SafeName": "Notification Engine"
}
{
  "SafeName": "PVWAReports"
}
{
  "SafeName": "PVWATicketingSystem"
}
{
  "SafeName": "PVWAPublicData"
}
{
  "SafeName": "PasswordManager"
}
{
  "SafeName": "PasswordManager_Pending"
}
{
  "SafeName": "AccountsFeedADAccounts"
}
{
  "SafeName": "AccountsFeedDiscoveryLogs"
}
```

### Application

Full example available at [dev/main.go]().

#### Import into project

`github.com/infamousjoeg/pas-api-go/pkg/cybr/api`

```go
package main

import pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
```

#### Logon to the PAS REST API Web Service

```go
package main

import (
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
)

var (
	hostname = os.Getenv("PAS_HOSTNAME")
	username = os.Getenv("PAS_USERNAME")
	password = os.Getenv("PAS_PASSWORD")
	authType = os.Getenv("PAS_AUTH_TYPE")
)

func main() {
	// Logon to PAS REST API Web Services
	token, errLogon := pasapi.Logon(hostname, username, password, authType, false)
	if errLogon != nil {
		log.Fatalf("Authentication failed. %s", errLogon)
	}
	fmt.Printf("Session Token:\r\n%s\r\n\r\n", token)
```

#### Call functions by referencing `pasapi` and "dot-referencing"

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	pasapi "github.com/infamousjoeg/pas-api-go/pkg/cybr/api"
)

// Declare variables (using Summon so they are env vars)
var (
	hostname = os.Getenv("PAS_HOSTNAME")
	username = os.Getenv("PAS_USERNAME")
	password = os.Getenv("PAS_PASSWORD")
	authType = os.Getenv("PAS_AUTH_TYPE")
)

// Start main function
func main() {
   // Verify PAS REST API Web Services
   // --> pasapi.ServerVerify is a "dot-reference"
	resVerify := pasapi.ServerVerify(hostname)
   // Marshal (convert) struct data into JSON format for pretty print
   // Otherwise, we "dot-reference" e.g. jsonVerify.ApplicationName would equal 'PasswordVault'
	jsonVerify, err := json.Marshal(&resVerify)
	if err != nil {
		log.Fatalf("Unable to marshal struct to JSON for verify. %s", err)
	}
   fmt.Printf("Verify JSON:\r\n%s\r\n\r\n", jsonVerify)
}
```

## Required Environment Variables

| Variable Name | Description |
| --- | --- |
| `PAS_HOSTNAME` | Base URL for PAS REST API Web Service |
| `PAS_USERNAME` | Username to authn to PAS REST API |
| `PAS_PASSWORD` | Password associated with `PAS_USERNAME` |
| `PAS_AUTH_TYPE` | Authentication method to use (cyberark or ldap) |

## Testing

1. Download and install [summon](https://cyberark.github.io/summon).
2. [OPTIONAL] Download and install a [summon provider](https://cyberark.github.io/summon/#providers).
   1. I use the `keyring` provider with [conceal](https://github.com/infamousjoeg/conceal).
3. Modify the values in [secrets.yml]() for your environment.
   1. If you did not complete the optional Step #2, you will use literal strings for `PAS_USERNAME` and `PAS_PASSWORD` similar to the values of `PAS_HOSTNAME` and `PAS_AUTH_TYPE`.
4. Run [main.go]() with the command: `summon go run main.go`.

### Successful Output

```shell
$ summon go run dev/main.go

Verify JSON:
{"ApplicationName":"PasswordVault","AuthenticationMethods":[{"Enabled":false,"Id":"windows"},{"Enabled":false,"Id":"pki"},{"Enabled":true,"Id":"cyberark"},{"Enabled":false,"Id":"oraclesso"},{"Enabled":false,"Id":"rsa"},{"Enabled":true,"Id":"radius"},{"Enabled":true,"Id":"ldap"},{"Enabled":true,"Id":"saml"}],"ServerId":"e00e8q16-b637-11e9-8329-ccd02f0167674","ServerName":"Vault"}

Session Token:
ZDRjNjNjNGItMVPjMS00MzRhLWIyNWMtYzI3MjMzZDFjNDg0OzNCRjk4NDEyMjEyNzgyOUI7MDAwMDAwMDJFN0ZERDcyNzJENUM3MkNDRjdBNUNDQ0UxQjY1QTYyMTkyMTlDQ0I0NTdGMjgxNDkxOTc1RTQxMjc1MkRFRTRFMDAwMDAwMDA7

List Safes JSON:
{"Safes":[{"SafeUrlId":"VaultInternal","SafeName":"VaultInternal","Location":"\\"},{"SafeUrlId":"Notification%20Engine","SafeName":"Notification Engine","Location":"\\"},{"SafeUrlId":"PVWAReports","SafeName":"PVWAReports","Location":"\\"},{"SafeUrlId":"PVWATicketingSystem","SafeName":"PVWATicketingSystem","Location":"\\"},{"SafeUrlId":"PVWAPublicData","SafeName":"PVWAPublicData","Location":"\\"},{"SafeUrlId":"PasswordManager","SafeName":"PasswordManager","Location":"\\"},{"SafeUrlId":"PasswordManager_Pending","SafeName":"PasswordManager_Pending","Location":"\\"},{"SafeUrlId":"AccountsFeedADAccounts","SafeName":"AccountsFeedADAccounts","Location":"\\"},{"SafeUrlId":"AccountsFeedDiscoveryLogs","SafeName":"AccountsFeedDiscoveryLogs","Location":"\\"},{"SafeUrlId":"D-Nix-Root","SafeName":"D-Nix-Root","Location":"\\"},{"SafeUrlId":"D-Win-DomainAdmins","SafeName":"D-Win-DomainAdmins","Location":"\\"},{"SafeUrlId":"D-Win-LocalAdmins","SafeName":"D-Win-LocalAdmins","Location":"\\"},{"SafeUrlId":"D-AWS-AccessKeys","SafeName":"D-AWS-AccessKeys","Location":"\\"},{"SafeUrlId":"D-Nix-AWSEC2-Keypairs","SafeName":"D-Nix-AWSEC2-Keypairs","Location":"\\"},{"SafeUrlId":"D-App-CyberArk-API","SafeName":"D-App-CyberArk-API","Location":"\\"},{"SafeUrlId":"D-MySQL-Users","SafeName":"D-MySQL-Users","Location":"\\"},{"SafeUrlId":"D-Nix-AWS-EC2","SafeName":"D-Nix-AWS-EC2","Location":"\\"},{"SafeUrlId":"PSMUniversalConnectors","SafeName":"PSMUniversalConnectors","Location":"\\"},{"SafeUrlId":"D-App-Docker-Registry","SafeName":"D-App-Docker-Registry","Location":"\\"},{"SafeUrlId":"D-Win-SvcAccts","SafeName":"D-Win-SvcAccts","Location":"\\"},{"SafeUrlId":"D-Nix-Root-2","SafeName":"D-Nix-Root-2","Description":"REST API Summit Toronto","Location":"\\"},{"SafeUrlId":"DemoSafe","SafeName":"DemoSafe","Location":"\\"},{"SafeUrlId":"Version1Safe","SafeName":"Version1Safe","Location":"\\"},{"SafeUrlId":"DemoSafe1","SafeName":"DemoSafe1","Location":"\\"},{"SafeUrlId":"Dummy-Safe","SafeName":"Dummy-Safe","Description":"Dummy safe","Location":"\\"}]}

List Safe Members JSON:
{"members":[{"Permissions":{"Add":true,"AddRenameFolder":true,"BackupSafe":true,"Delete":true,"DeleteFolder":true,"ListContent":true,"ManageSafe":true,"ManageSafeMembers":true,"MoveFilesAndFolders":true,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":true,"ViewAudit":true,"ViewMembers":true},"UserName":"jgarcia"},{"Permissions":{"Add":true,"AddRenameFolder":true,"BackupSafe":true,"Delete":true,"DeleteFolder":true,"ListContent":true,"ManageSafe":true,"ManageSafeMembers":true,"MoveFilesAndFolders":true,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":true,"ViewAudit":true,"ViewMembers":true},"UserName":"Master"},{"Permissions":{"Add":true,"AddRenameFolder":true,"BackupSafe":true,"Delete":true,"DeleteFolder":true,"ListContent":true,"ManageSafe":true,"ManageSafeMembers":true,"MoveFilesAndFolders":true,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":true,"ViewAudit":true,"ViewMembers":true},"UserName":"Batch"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":true,"Delete":false,"DeleteFolder":false,"ListContent":false,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":false,"ViewMembers":false},"UserName":"Backup Users"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"Auditors"},{"Permissions":{"Add":false,"AddRenameFolder":true,"BackupSafe":false,"Delete":false,"DeleteFolder":true,"ListContent":false,"ManageSafe":true,"ManageSafeMembers":false,"MoveFilesAndFolders":true,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":true,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":false,"ViewMembers":false},"UserName":"Operators"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":true,"Delete":false,"DeleteFolder":false,"ListContent":false,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":false,"ViewMembers":false},"UserName":"DR Users"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"Notification Engines"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"PVWAGWAccounts"},{"Permissions":{"Add":true,"AddRenameFolder":true,"BackupSafe":false,"Delete":true,"DeleteFolder":true,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":true,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":false},"UserName":"PasswordManager"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":true,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":false},"UserName":"D-Win-DomainAdmin_Users"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":false,"Retrieve":false,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"D-Win-DomainAdmin_Auditors"},{"Permissions":{"Add":true,"AddRenameFolder":false,"BackupSafe":true,"Delete":true,"DeleteFolder":false,"ListContent":true,"ManageSafe":true,"ManageSafeMembers":true,"MoveFilesAndFolders":false,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"D-Win-DomainAdmin_Admins"},{"Permissions":{"Add":true,"AddRenameFolder":true,"BackupSafe":true,"Delete":true,"DeleteFolder":true,"ListContent":true,"ManageSafe":true,"ManageSafeMembers":true,"MoveFilesAndFolders":true,"Rename":true,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":true,"Update":true,"UpdateMetadata":true,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"Vault Admins"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"AIMWebService"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"SlackBot"},{"Permissions":{"Add":false,"AddRenameFolder":false,"BackupSafe":false,"Delete":false,"DeleteFolder":false,"ListContent":true,"ManageSafe":false,"ManageSafeMembers":false,"MoveFilesAndFolders":false,"Rename":false,"RestrictedRetrieve":true,"Retrieve":true,"Unlock":false,"Update":false,"UpdateMetadata":false,"ValidateSafeContent":false,"ViewAudit":true,"ViewMembers":true},"UserName":"Prov_PASAAS-PVWA"}]}

List Applications JSON:
{"application":[{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"DemoApp","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"AIMWebService","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"SlackBot","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"Ansible","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"AD Automation","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"DockerRegistry","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"AnsibleCP","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"pyAIM","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"},{"AccessPermittedFrom":0,"AccessPermittedTo":24,"AllowExtendedAuthenticationRestrict":false,"AppID":"IDaptive","BusinessOwnerEmail":"","BusinessOwnerFName":"","BusinessOwnerLName":"","BusinessOwnerPhone":"","Description":"","Disabled":false,"ExpirationDate":"0001-01-01T00:00:00Z","Location":"\\Applications"}]}

List Application Authentication Methods JSON:
{"authentication":[{"AppID":"Ansible","AuthType":"certificateSerialNumber","AuthValue":"5e00000044f3868cc4ce21134c000000000004","authID":0},{"AppID":"Ansible","AuthType":"machineAddress","AuthValue":"12.233.234.444","authID":0},{"AppID":"Ansible","AuthType":"machineAddress","AuthValue":"54.142.78.106","authID":0},{"AppID":"Ansible","AuthType":"machineAddress","AuthValue":"54.235.118.11","authID":0},{"AppID":"Ansible","AuthType":"certificateattr","AuthValue":"","authID":0}]}

Successfully logged off PAS REST API Web Services.
```

# cybr-cli <!-- omit in toc -->

A "Swiss Army Knife" command-line interface (CLI) for easy human and non-human interaction with CyberArk's suite of products.

Current products supported:
* CyberArk Privileged Access Security (PAS)
  * Accounts, Safes, Safe Members, Applications
* CyberArk Enterprise Secrets Manager & Conjur
  * List Resources, Policy, Secrets, API Key Rotation, Authenticators

**Want to get dangerous quickly?** Check out the example bash script at [dev/add-delete-pas-application.sh]().

[![cybr-cli CI](https://github.com/infamousjoeg/cybr-cli/workflows/cybr-cli%20CI/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3A%22cybr-cli+CI%22) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=infamousjoeg_pas-api-go&metric=alert_status)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ALint) [![CodeQL](https://github.com/infamousjoeg/cybr-cli/workflows/CodeQL/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ACodeQL) [![](https://img.shields.io/github/downloads/infamousjoeg/cybr-cli/latest/total?color=blue&label=Download%20Latest%20Release&logo=github)](https://github.com/infamousjoeg/cybr-cli/releases/latest)

## Table of Contents <!-- omit in toc -->

- [Install](#install)
  - [MacOS](#macos)
  - [Windows & Linux](#windows-or-linux)
- [Usage](#usage)
  - [Command-Line Interface (CLI)](#command-line-interface-cli)
    - [logon](#logon)
    - [logoff](#logoff)
    - [accounts](#accounts)
      - [list](#list)
      - [get](#get)
      - [add](#add)
      - [delete](#delete)
    - [safes](#safes)
      - [list](#list-1)
      - [add](#add-1)
      - [update](#update)
      - [delete](#delete-1)
      - [list-members](#list-members)
      - [add-member](#add-member)
      - [remove-member](#remove-member)
    - [applications](#applications)
      - [list](#list-2)
      - [add](#add-2)
      - [delete](#delete-2)
      - [list-authn](#list-authn)
      - [add-authn](#add-authn)
      - [delete-authn](#delete-authn)
    - [conjur](#conjur)
      - [logon](#logon-1)
      - [logon-non-interactive](#logon-non-interactive)
      - [rotate-api-key](#rotate-api-key)
      - [list](#list-3)
      - [append-policy](#append-policy)
      - [update-policy](#update-policy)
      - [replace-policy](#replace-policy)
      - [get-secret](#get-secret)
      - [set-secret](#set-secret)
      - [enable-authn](#enable-authn)
      - [info](#info)
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
      - [Environment Variables Used](#environment-variables-used)
- [Testing](#testing)
- [Maintainers](#maintainers)
- [Contributions](#contributions)
- [License](#license)

## Install

### MacOS

```shell
$ brew tap infamousjoeg/tap
$ brew install cybr-cli
```

### Windows or Linux

Download from the [Releases](https://github.com/infamousjoeg/cybr-cli/releases) page.

## Usage

### Command-Line Interface (CLI)

#### logon

```shell
$ cybr logon -u username -a cyberark-or-ldap -b https://pvwa.example.com
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-u|--username|☑||Logon username|jgarcia|
|-a|--auth-type|☑||Authentication method|ldap|
|-b|--base-url|☑||URL to /PasswordVault|https://pvwa.example.com|
|-i|--insecure-tls||false|Whether to validate TLS|false|

Logon to the PAS REST API as the username you provide using the authentication method you choose. At this time, only `cyberark` and `ldap` authentication methods are supported.

Upon successful logon, a file will be created in your user's home directory at `.cybr/config`. It is an encoded file that cannot be read in plain-text. This holds your current session information.

#### logoff

```shell
$ cybr logoff
```

Logoff the PAS REST API as the username you provided during logon.

Upon successful logoff, the config file located in your user's home directory at `.cybr/config` will be removed and the session token stored within will be expired.

#### accounts

```shell
$ cybr accounts
```

Display help for the `cybr accounts` command.

##### list

```shell
$ cybr accounts list
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--search|||List of keywords to search for in accounts, separated by space|"domain windows Test-Safe"|
|-t|--search-type||contains|Get accounts that either contain or start with the value specified in the Search parameter. Valid values: contains (default) or startswith|startswith|
|-r|--sort|||Property or properties by which to sort returned accounts, followed by asc (default) or desc to control sort direction. Separate multiple properties with commas, up to a maximum of three (3) properties|"name,address,port desc"|
|-o|--offset||0|Offset of the first account that is returned in the collection of results|51|
|-l|--limit||50|Maximum number of returned accounts. If not specified, the default value is 50. The maximum number that can be specified is 1000|200|
|-f|--filter|||Search accounts filtered by safename or modificationTime|"safename eq Test-Safe"|

List all accounts the logged on user can read.

##### get

```shell
$ cybr accounts get -i 24_1
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-i|--account-id|☑||Account ID (not name) to get account details for|24_1|

_**Note:** The Account ID is a unique primary key within the Secure Digital Vault for CyberArk PAS. This is NOT equivalent to the object `name` or account `username` properties._

Get account object details based on Account ID specified for the account.

##### add

```shell
$ cybr accounts add -s SafeName -p PlatformID -u Username -a 10.0.0.1 -t password -s SuperSecret
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-n|--name|||The name of the account object being created. Will use auto-generated name if not provided|DeviceType-PlatformName-SafeName-Username|
|-a|--address|☑||Address of the account object|10.0.0.1|
|-u|--username|☑||Username of the account object|root|
|-p|--platform-id|☑||Platform ID of the account object|WinDomain|
|-s|--safe|☑||Safe name to store the account object within|Test-Safe|
|-t|--secret-type|☑||Secret type of the account object. e.g. password, accessKey, sshKey|password|
|-c|--secret|☑||Secret of the account object|SuperSecret|
|-e|--platform-properties|||Extra platform properties|port=22,UseSudoOnReconcile=yes,CustomField=customValue|
|-m|--automatic-management||false|If set, will automatically manage the onboarded account|true|
|-r|--manual-management-reason|||The reason the account object is not being managed|"Testing onboarding"|

Add an account object to a safe.

##### delete

```shell
$ cybr accounts delete -i 24_1
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-i|--account-id|☑||Account ID (not name) to delete|24_1|

_**Note:** The Account ID is a unique primary key within the Secure Digital Vault for CyberArk PAS. This is NOT equivalent to the object `name` or account `username` properties._

Delete a specific account object from within a safe. The account will be marked for deletion until the safe's retention policy period has expired.

#### safes

```shell
$ cybr safes
```

Display help for the `cybr safes` command.

##### list

```shell
$ cybr safes list
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|||Safe name to target specifically|P-WIN-ADMINS-DOMAIN|

List all safes the username you are logged on as has access to read. If the `-s` or `--safe` optional flag is given, only that targeted safe's details will be returned.

##### add

```shell
$ cybr safes add -s SafeName -d Description --cpm ManagingCPM --days 0
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|☑||Safe name to create|P-WIN-ADMINS-DOMAIN|
|-d|--desc|☑||Description of the safe created||
|-O|--olac||false|Enable object-level access control|false|
||--cpm||PasswordManager|Set the Managing CPM user|PasswordManager1|
||--days||7|Number of days to retain password versions for|0|
|-P|--auto-purge||false|This should not be needed|false|
|-l|--location||\\|The location of the Safe in the Secure Digital Vault|\\|

Add a safe and configure it's retention and location.

##### update

```shell
$ cybr safes update -t TargetSafeName -s NewSafeName -d NewDesc
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-t|--target-safe|☑||Safe name to update|
|-s|--safe|||New name of safe|P-WIN-ADMINS-DOMAIN|
|-d|--desc|||New description of safe||
|-O|--olac||false|Enable object-level access control on safe (this is not reversible)|false|
||--cpm|||New managing CPM user to change to|PasswordManager2|

Update a safe. Only the options provided will be modified.

##### delete

```shell
$ cybr safes delete -s SafeName
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|☑||Safe name to delete|P-WIN-ADMINS-DOMAIN|

Delete a safe. If the safe has a retention policy of days that is greater than 0, the safe will be marked for deletion until the amount of days assigned are met.

##### list-members

```shell
$ cybr safes list-members -s SafeName
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|☑||Safe name to list members from|P-WIN-ADMINS-DOMAIN|

List all safe members, or a specific safe, that the user logged on can read.

##### add-member

```shell
$ cybr safes add-member -s SafeName -m MemberName --retrieve-account
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|☑||Safe name to add member to|P-WIN-ADMINS-DOMAIN|
|-m|--member-name|☑||Name of user or group member being added to safe|"Vault Admins"|
|-i|--search-in||Vault|Search in Domain or `Vault` for user or group member to add|ldap-domain.com|
|-e|--member-expiration-date|||When the membership will expire|2022-01-01|
||--use-accounts||false|Use accounts in safe|true|
||--retrieve-accounts||false|Retrieve accounts in safe|true|
||--update-account-content||false|Update account content in safe|true|
||--update-account-properties||false|Update account properties in safe|true|
||--init-cpm-account-management-operations||false|Perform automated CPM actions on accounts within the safe|true|
||--specify-next-account-content||false|Specify next account's content within safe|true|
||--manage-safe||false|Manage the safe|true|
||--manage-safe-members||false|Manage members of the safe|true|
||--backup-safe||false|Backup the safe|true|
||--view-audit-log||false|View audit log of safe|true|
||--view-safe-members||false|View the safe members|true|
||--access-content-without-confirmation||false|Access the Vault without needing approvals|true|
||--create-folders||false|Create folders within safe|true|
||--delete-folders||false|Delete folders within safe|true|
||--move-accounts-and-folders||false|Move accounts and folders to other safes|true|

Adds an existing user as a Safe member. The user who runs this command required `Manage Safe Members` permission on the safe.

##### remove-member

```shell
$ cybr safes remove-member -s SafeName -m MemberName
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--safe|☑||Safe name to delete member from|P-WIN-ADMINS-DOMAIN|
|-m|--member-name|☑||Name of user or group member being added to safe|"Vault Admins"|

Removes a specific member from a Safe. The user who runs this command requires `Manage Safe Members` permission on the safe.

#### applications

```shell
$ cybr applications
```

Display help for the `cybr applications` command.

##### list

```shell
$ cybr applications list
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-l|--location||\\|Location of application in EPV|\\|

List all applications the username you are logged on as has access to read.

##### add

```shell
$ cybr applications add -a AppID -l "\\"
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|
|-l|--location|☑||Application location|\\|
|-d|--description||Application description||
|-f|--access-permitted-from||0|Access permitted for the application (0-23)|0|
|-t|--access-permitted-to||23|Access permitted for the application (0-23)|23|
|-e|--expiration-date|||When application will expire||
|-i|--disabled|||Disable the application, e.g. yes/no|yes|
|-r|--business-owner-first-name|||Application business owner first name|Joe|
|-n|--business-owner-last-name|||Application business owner last name|Garcia|
|-m|--business-owner-email|||Application business owner email|email@example.com|
|-p|--business-owner-phone|||Application business owner phone|555-555-1234|

Add an application identity.

##### delete

```shell
$ cybr applications delete -a AppID
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|

Delete an application identity.

##### list-authn

```shell
$ cybr applications list-authn -a AppID
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|

List all authentication methods configured for the application identity given.

##### add-authn

```shell
$ cybr applications add-authn -a AppID -t path -v /some/path
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|
|-t|--auth-type|☑||Application authentication method type|allowedMachines|
|-v|--auth-value|☑||Application authentication method type value|10.0.0.1|
|-f|--is-folder||false|Application is a folder|true|
|-s|--allow-internal-scripts||false|Allow internal scripts|true|

Add an authentication method to an application identity.

##### delete-authn

```shell
$ cybr applications delete-authn -a AppID -i 1
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|
|-i|--auth-method-id|☑|Application authentication method ID to be deleted|1|

Delete an authentication method of an application identity.

#### conjur

```shell
$ cybr conjur
```
Display help for the `cybr conjur` command.

##### logon

```shell
$ cybr conjur logon -a account -b https://conjur.example.com -l admin
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-l|--login|☑||Conjur login name|admin|
|-a|--account|☑||Conjur organization account name|cyberark|
|-b|--base-url|☑||Conjur appliance URL|https://conjur.example.com|

Authenticate to Conjur using API Key or Password.

##### logon-non-interactive

```shell
$ cybr conjur logon-non-interactive
```

Authenticate to Conjur using API Key or Password non-interactively. This requires `.netrc` and `.conjurrc` files to be located in the User's home directory. These will be used for authentication purposes.

##### rotate-api-key

```shell
$ cybr conjur rotate-api-key -l host/some/application
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-l|--login|||The username of the resource to rotate the API key of|host/ansible/tower|

Replaces the API Key of another role you can update with a new, securely random API Key. The new API Key is returned to the STDOUT pipe.

##### list

```shell
$ cybr conjur list --kind variable --search prod
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-k|--kind|||Filters results to only resources of that kind|hosts|
|-s|--search|||Filters results to those pertaining to the search query|aws|
|-l|--limit||10|Maximum number of returned resources. Default is ten (10)|25|
|-o|--offset||0|Index to start returning results from for pagination|26|
|-i|--inspect||false|Show full object information|true|

Lists resources within an organization account.

##### append-policy

```shell
$ cybr conjur append-policy --branch root --file ./path/to/root.yml
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-b|--branch|☑||The policy branch in which policy is being loaded|root|
|-f|--file|☑||The policy file that will be loaded into the branch|./path/to/root.yml|

Adds data to the existing Conjur policy. Deletions are not allowed. Any policy objects that exist on the server but are omitted from the policy file will not be deleted and any explicit deletions in the policy file will result in an error.

##### update-policy

```shell
$ cybr conjur update-policy --branch root --file ./path/to/root.yml
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-b|--branch|☑||The policy branch in which policy is being loaded|root|
|-f|--file|☑||The policy file that will be loaded into the branch|./path/to/root.yml|

Modifies an existing Conjur policy. Data may be explicitly deleted using the `!delete`, `!revoke`, and `!deny` statements. Unlike `replace` mode, no data is ever implicitly deleted.

##### replace-policy

```shell
$ cybr conjur replace-policy --branch root --file ./path/to/root.yml
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-b|--branch|☑||The policy branch in which policy is being loaded|root|
|-f|--file|☑||The policy file that will be loaded into the branch|./path/to/root.yml|

Loads or replaces a Conjur policy. Any policy data which already exists on the server but is not explicitly specified in the new policy file **will be deleted**.

##### get-secret

```shell
$ cybr conjur get-secret -i id/to/variable
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-i|--id|☑||The variable ID containing the secret|/aws/us-east-1/access_key_id|
|-n|--no-new-line||false|Return the value without a new line|true|

Fetches the value of a secret from the specified Variable.

##### set-secret

```shell
$ cybr conjur set-secret -i id/to-variable -v "P@$$word"
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-i|--id|☑||The variable ID being updated|/aws/us-east-1/secret_access_key|
|-n|--no-new-line||false|Return the value without a new line|true|

Sets a secret value for the specified Variable.

##### enable-authn

```shell
$ cybr conjur enable-authn -s authn-iam/prod
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-s|--service-id|☑||The authenticator service ID|authn-k8s/k8s-cluster-1|

Enables a Conjur authenticator.

##### info

```shell
$ cybr conjur info
```

|Short|Long|Required|Default Value|Description|Example|
|---|---|---|---|---|---|
|-a|--app-id|☑||Application ID|Ansible|

Get info about the Conjur cluster and health.

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
cybr v0.0.5-alpha
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
	hostname = os.Getenv("PAS_BASE_URL")
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
	hostname = os.Getenv("PAS_BASE_URL")
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

##### Environment Variables Used

| Variable Name | Description |
| --- | --- |
| `PAS_BASE_URL` | Base URL for PAS REST API Web Service |
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

## Maintainers

[@infamousjoeg](https://github.com/infamousjoeg)

[![Buy me a coffee][buymeacoffee-shield]][buymeacoffee]

[buymeacoffee]: https://www.buymeacoffee.com/infamousjoeg
[buymeacoffee-shield]: https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png

[@AndrewCopeland](https://github.com/AndrewCopeland)

## Contributions

Pull Requests are currently being accepted.  Please read and follow the guidelines laid out in [CONTRIBUTING.md]().

## License

[Apache 2.0](LICENSE)

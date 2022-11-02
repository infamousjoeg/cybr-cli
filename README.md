# cybr-cli <!-- omit in toc -->

A "Swiss Army Knife" command-line interface (CLI) for easy human and non-human interaction with CyberArk's suite of products.

Current products supported:
* CyberArk Privileged Access Manager (PAM)
* CyberArk Secrets Manager Central Credential Provider (CCP)
* CyberArk Conjur Secrets Manager Enterprise & [Open Source](https://conjur.org)
* CyberArk Cloud Entitlements Manager ([Free trial](https://www.cyberark.com/try-buy/cloud-entitlements-manager/))

**Want to get dangerous quickly?** Check out the example bash script at [dev/add-delete-pas-application.sh](dev/add-delete-pas-application.sh).

[![cybr-cli CI](https://github.com/infamousjoeg/cybr-cli/workflows/cybr-cli%20CI/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3A%22cybr-cli+CI%22) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=infamousjoeg_pas-api-go&metric=alert_status)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ALint) [![CodeQL](https://github.com/infamousjoeg/cybr-cli/workflows/CodeQL/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ACodeQL) [![](https://img.shields.io/github/downloads/infamousjoeg/cybr-cli/latest/total?color=blue&label=Download%20Latest%20Release&logo=github)](https://github.com/infamousjoeg/cybr-cli/releases/latest)

## Table of Contents <!-- omit in toc -->

- [Install](#install)
	- [MacOS](#macos)
	- [Windows or Linux](#windows-or-linux)
	- [AWS CloudShell](#aws-cloudshell)
	- [Install from Source](#install-from-source)
- [Usage](#usage)
	- [Authenticating with authn-iam (AWS IAM Role Authentication)](#authenticating-with-authn-iam-aws-iam-role-authentication)
	- [Documentation](#documentation)
- [Autocomplete](#autocomplete)
- [Example Source Code](#example-source-code)
	- [Logon to the PAS REST API Web Service](#logon-to-the-pas-rest-api-web-service)
- [Security](#security)
	- [`cybr safes add-member --role` Role Permissions](#cybr-safes-add-member---role-role-permissions)
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


### AWS CloudShell

```shell
mkdir -p ~/.local/bin && \
curl --silent "https://api.github.com/repos/infamousjoeg/cybr-cli/releases/latest" |
    grep '"tag_name":' |
    sed -E 's/.*"([^"]+)".*/\1/' |
    xargs -I {}  curl -o ~/.local/bin/cybr -sOL "https://github.com/infamousjoeg/cybr-cli/releases/download/"{}'/linux_cybr' && \
chmod +x ~/.local/bin/cybr
```

### Install from Source

```shell
$ git clone https://github.com/infamousjoeg/pas-api-go.git
$ ./install
$ cybr help
```

## Usage

* `$ cybr help` for top-level commands list
* `$ cybr [command] -h` for specific command details and sub-commands list

### Authenticating with authn-iam (AWS IAM Role Authentication)

Set the following environment variables:

* `CONJUR_ACCOUNT` - The Conjur account name
* `CONJUR_APPLIANCE_URL` - The URL of the Conjur service (e.g. https://conjur.example.com)
* `CONJUR_AUTHN_LOGIN` - The Host ID for the IAM role (e.g. `host/cloud/aws/ec2/1234567890/ConjurAWSRoleEC2`)
* `CONJUR_AUTHENTICATOR` - The authenticator ID (e.g. `authn-iam`)
* `CONJUR_AUTHN_SERVICE_ID` - The authenticator web service ID (e.g. `prod`)
* `CONJUR_AWS_TYPE` - The AWS type (e.g. `ec2` or `ecs` or `lambda`)

Once environment variables are set, ensure no .conjurrc or .netrc exists in the user's home directory:

`rm -f ~/.conjurrc ~/.netrc`

Then run any command you wish to run within `cybr conjur`. Use the `--help` flag to see all available commands.

### Documentation

All commands are documentated [in the docs/ directory](docs/cybr.md).

## Autocomplete
The `cybr` CLI has a `completion` command that can be used to enable autocomplete for the CLI.
The completion command is dependant on your shell type. Currently the only shells that are supported are: bash, zsh, fish and powershell.

Below is an example on how to enable `cybr` cli auto-completion from a zsh shell.
```bash
# enable shell completetion. Only needs to be performed once.
echo "autoload -U compinit; compinit" >> ~/.zshrc

# create and write the auto-completion script.
# ${fpath[1]} '1' may be different depending on your environment.
cybr completion zsh > "${fpath[1]}/_cybr"
```

If you are using a different shell execute the `completion` command with the `--help` flag and follow instructions for the desired shell type.
```bash
cybr completion --help
```

## Example Source Code

### Logon to the PAS REST API Web Service

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
}
```

## Security

If there is a security concern or bug discovered, please responsibly disclose all information to joe (dot) garcia (at) cyberark (dot) com.

### `cybr safes add-member --role` Role Permissions

All safe member roles defined below are based on best practices and recommendations put forth by CyberArk's PAS Programs Office, creators of the CyberArk Blueprint for Identity Security.

|Role|Safe Authorizations|
|---|---|
|BreakGlass|All authorizations except Authorize Password Requests|
|VaultAdmin|- List Accounts<br>- View Audit Log<br>- View Safe Members|
|SafeManager|- Manage Safe<br>- Manage Safe Members<br>- View Audit Log<br>- View Safe Members<br>- Access Safe w/o Confirmation|
|EndUser|- Use/Retrieve/List Accounts<br>- View Audit Log<br>- View Safe Members|
|Auditor|- List Accounts<br>- View Audit Log<br>- View Safe Members|
|AIMWebService|No authorizations|
|AppProvider|- Retrieve/List Accounts<br>- View Safe Members|
|ApplicationIdentity|- Retrieve/List Accounts|
|AccountProvisioner|- List/Add/Delete Accounts<br>- Update Password Properties<br>- Initiate CPM Password Management Operations<br>- View Audit Log<br>- View Safe Members<br>- Access Safe w/o Confirmation|
|CPDeployer|- List/Add Accounts<br>- Update Password Properties<br>- Initiate CPM Password Management Operations<br>- Manage Safe Member<br>- View Audit Log, View Safe Members<br>- Access Safe w/o Confirmation|
|ComponentOrchestrator|- List/Add Accounts<br>- Update Password Properties<br>- Initiate CPM Password Management Operations<br>- View Audit Log<br>- Access Safe w/o Confirmation|
|APIAutomation|- List/Add/Rename/Delete/Unlock Accounts<br>- Update Password Content/Properties<br>- Initiate CPM Password Management Operations<br>- Manage Safe<br>- Manage Safe Members<br>- View Audit Log<br>- View Safe Members<br>- Create/Delete Folders<br>- Move Accounts/Folders|
|PasswordScheduler|- List Accounts<br>- Initiate CPM Password Management Operation<br>- View Audit Log<br>- View Safe Members<br>- Access Safe w/o Confirmation|
|ApproverLevel1|- List Accounts<br>- View Audit Log<br>- View Safe Members<br>- Authorize Password Requests (Level 1)|
|ApproverLevel2|- List Acccounts<br>- View Audit Log<br>- View Safe Members<br>- Authorize Password Requests (Level 2)|

## Testing

`go test -v ./...`

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

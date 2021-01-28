# cybr-cli <!-- omit in toc -->

A "Swiss Army Knife" command-line interface (CLI) for easy human and non-human interaction with CyberArk's suite of products.

Current products supported:
* CyberArk Privileged Access Security (PAS)
* CyberArk Conjur Secrets Manager Enterprise & Open Source

**Want to get dangerous quickly?** Check out the example bash script at [dev/add-delete-pas-application.sh]().

[![cybr-cli CI](https://github.com/infamousjoeg/cybr-cli/workflows/cybr-cli%20CI/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3A%22cybr-cli+CI%22) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=infamousjoeg_pas-api-go&metric=alert_status)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ALint) [![CodeQL](https://github.com/infamousjoeg/cybr-cli/workflows/CodeQL/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ACodeQL) [![](https://img.shields.io/github/downloads/infamousjoeg/cybr-cli/latest/total?color=blue&label=Download%20Latest%20Release&logo=github)](https://github.com/infamousjoeg/cybr-cli/releases/latest)

## Table of Contents <!-- omit in toc -->

- [Install](#install)
  - [MacOS](#macos)
  - [Windows or Linux](#windows-or-linux)
  - [Install from Source](#install-from-source)
  - [Docker Container](#docker-container)
    - [Run Container Indefinitely](#run-container-indefinitely)
    - [Run Container Ephemerally (Recommended)](#run-container-ephemerally-recommended)
      - [One-Time Use](#one-time-use)
      - [One-Time Use w/ Saved Config](#one-time-use-w-saved-config)
      - [Using with jq](#using-with-jq)
- [Usage](#usage)
- [Example Source Code](#example-source-code)
  - [Logon to the PAS REST API Web Service](#logon-to-the-pas-rest-api-web-service)
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
cybr v0.1.0-beta
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

## Usage

* `$ cybr help` for top-level commands list
* `$ cybr [command] -h` for specific command details and sub-commands list

All commands are documentated [in the docs/ directory](docs/cybr.md).

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
```

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

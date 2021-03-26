# cybr-cli <!-- omit in toc -->

A "Swiss Army Knife" command-line interface (CLI) for easy human and non-human interaction with CyberArk's suite of products.

Current products supported:
* CyberArk Privileged Access Security (PAS)
* CyberArk Conjur Secrets Manager Enterprise & Open Source

**Want to get dangerous quickly?** Check out the example bash script at [dev/add-delete-pas-application.sh](dev/add-delete-pas-application.sh).

[![cybr-cli CI](https://github.com/infamousjoeg/cybr-cli/workflows/cybr-cli%20CI/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3A%22cybr-cli+CI%22) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=infamousjoeg_pas-api-go&metric=alert_status)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ALint) [![CodeQL](https://github.com/infamousjoeg/cybr-cli/workflows/CodeQL/badge.svg)](https://github.com/infamousjoeg/cybr-cli/actions?query=workflow%3ACodeQL) [![](https://img.shields.io/github/downloads/infamousjoeg/cybr-cli/latest/total?color=blue&label=Download%20Latest%20Release&logo=github)](https://github.com/infamousjoeg/cybr-cli/releases/latest)

## Table of Contents <!-- omit in toc -->

- [Install](#install)
  - [MacOS](#macos)
  - [Windows or Linux](#windows-or-linux)
  - [AWS CloudShell](#aws-cloudshell)
  - [Install from Source](#install-from-source)
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

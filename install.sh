#!/bin/bash
set -eou pipefail

wget https://github.com/infamousjoeg/cybr-cli/releases/latest/download/cybr-cli_linux_amd64.tar.gz
sudo tar -xzvf cybr-cli_linux_amd64.tar.gz -C /usr/local/bin

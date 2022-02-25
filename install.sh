#!/bin/bash
set -eou pipefail
IFS=$'/t/n'

wget https://github.com/infamousjoeg/cybr-cli/releases/latest/download/linux_cybr
sudo chmod +x linux_cybr
sudo mv linux_cybr /usr/bin/cybr

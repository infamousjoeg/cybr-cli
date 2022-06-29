#!/bin/bash

echo "Making Linux x64 binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux/amd64/cybr .
echo "Compressing Linux x64 binary into .tar.gz..."
tar -czf ./bin/linux/amd64/cybr-cli_linux_amd64.tar.gz ./bin/linux/amd64/cybr
echo "Making Linux ARM binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bin/linux/arm64/cybr .
echo "Compressing Linux ARM binary into .tar.gz..."
tar -czf ./bin/linux/arm64/cybr-cli_linux_arm64.tar.gz ./bin/linux/arm64/cybr
echo "Making Darwin x64 binary..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin/amd64/cybr .
echo "Requires notarization with Apple... [$PWD/bin/notarize.sh]"
echo "Making Darwin ARM binary..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin/arm64/cybr .
echo "Requires notarization with Apple... [$PWD/bin/notarize.sh]"
echo "Making Windows x64 binary..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/windows/amd64/cybr.exe .
echo "Compressing Windows x64 binary into .zip..."
zip -r ./bin/windows/amd64/cybr-cli_windows_amd64.zip ./bin/windows/amd64/cybr.exe
echo "Making Windows ARM binary..."
CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o ./bin/windows/arm64/cybr.exe .
echo "Compressing Windows ARM binary into .zip..."
zip -r ./bin/windows/arm64/cybr-cli_windows_arm64.zip ./bin/windows/arm64/cybr.exe
echo "Finished making - files output to directories in $PWD/bin"
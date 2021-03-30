#!/bin/bash

echo "Making Linux x64 binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux_cybr .
echo "Making Darwin x64 binary..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./bin/darwin_cybr .
echo "Making Darwin ARM binary..."
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o ./bin/darwin_arm64_cybr .
echo "Making Windows x64 binary..."
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./bin/windows_cybr.exe .
echo "Finished making - files output to ./bin/"
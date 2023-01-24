#!/bin/bash

# Array of all OSes to build for
os=("windows" "linux" "darwin")
# Array of all architectures to build for
arch=("amd64" "arm64")

# Loop through all OSes
for i in "${os[@]}"
do
    # Loop through all architectures
    for j in "${arch[@]}"
    do
        # Release the binary
        make release GOOS="$i" GOARCH="$j"
    done
done
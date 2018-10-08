#!/bin/bash

version=$1

echo "Creating build directory: $(pwd)/build/${version} ..."
mkdir -p build/$version

platforms=("windows/amd64" "linux/amd64" "darwin/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=nuxctl'-'$GOOS'-'$GOARCH

    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building nuxct for ${GOOS}/${GOARCH}..."
    env GOOS=$GOOS GOARCH=$GOARCH go build -o build/$version/$output_name main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
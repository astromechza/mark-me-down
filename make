#!/usr/bin/env bash

set -e

function buildbinary {
    goos=$1
    goarch=$2

    echo "Building official $goos $goarch binary"

    outputfolder="build/${goos}_${goarch}"
    echo "Output Folder $outputfolder"
    mkdir -pv $outputfolder

    export GOOS=$goos
    export GOARCH=$goarch

    go build -i -v -o "$outputfolder/mark-me-down" github.com/AstromechZA/mark-me-down

    echo "Done"
    ls -l "$outputfolder/mark-me-down"
    file "$outputfolder/mark-me-down"
    echo
}

# build for mac
buildbinary darwin amd64

# build for linux
buildbinary linux amd64

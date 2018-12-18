#!/bin/bash

if [ -d "release" ]; then
    rm -fr release
fi

mkdir release

build() {
    local os=$1
    local arch=$2
    local filename="credman-$os-$arch"
    local ext=""

    if [ "$os" == "windows" ]; then
        ext=".exe"
    fi

    echo "Compiling for $os-$arch"
    GOOS=$os GOARCH=$arch go build -o "release/$filename$ext" -ldflags='-s -w'

    echo "  Creating archive"
    tar -czf "release/$filename.tar.gz" "release/$filename$ext"

    echo "  Removing binary"
    rm -fr "release/$filename$ext"
}

build linux amd64
build linux 386
build linux arm
build linux arm64
build windows amd64
build windows 386
build darwin amd64
build darwin 386

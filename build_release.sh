#!/bin/bash

CWD=$(pwd)

if [ -d "release" ]; then
    rm -fr release
fi

mkdir release

build() {
    local os=$1
    local arch=$2
    local filename="credman-$os-$arch"
    local ext=""

    cd "$CWD/cmd/credman"

    if [ "$os" == "windows" ]; then
        ext=".exe"
    fi

    echo "Compiling for $os-$arch"
    GOOS=$os GOARCH=$arch go build -ldflags='-s -w'

    cd "$CWD/release"
    echo "  Creating archive"
    tar -czf "$filename.tar.gz" "$filename$ext"

    echo "  Removing binary"
    rm -fr "$filename$ext"
}

build linux amd64
build linux 386
build linux arm
build linux arm64
build windows amd64
build windows 386
build darwin amd64
build darwin 386

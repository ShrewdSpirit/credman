#!/bin/bash

CWD=$(pwd)

PACKAGE='github.com/ShrewdSpirit/credman'

VERSION='0.10.0'
COMMITHASH=$(git rev-parse HEAD)

VERINFO="-X $PACKAGE/data.Version=$VERSION -X $PACKAGE/data.GitCommit=$COMMITHASH"

cd gui/assets
npm run build
gassets -d .

if [[ $# -eq 1 && $1 == "install" ]]; then
    cd "$CWD/cmd/credman"
    go install -ldflags="$VERINFO"
else
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

        if [[ $os == "windows" ]]; then
            ext=".exe"
        fi

        echo "Compiling for $os-$arch"
        GOOS=$os GOARCH=$arch go build -ldflags="$VERINFO -s -w" -o "credman$ext"

        chmod +x "credman$ext"

        # cd "$CWD/release"
        echo "  Creating archive"
        tar -czf "$filename.tar.gz" "credman$ext"
        mv "$filename.tar.gz" "$CWD/release"

        echo "  Removing binary"
        rm -fr "credman$ext" "$filename.tar.gz"
    }

    build linux amd64
    build linux 386
    build linux arm
    build linux arm64
    build windows amd64
    build windows 386
    build darwin amd64
    build darwin 386
fi

echo "Done"

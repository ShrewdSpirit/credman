#!/bin/bash

CWD=$(pwd)

PACKAGE='github.com/ShrewdSpirit/credman'

VERSION='0.11.0'
COMMITHASH=$(git rev-parse HEAD)

VERINFO="-X $PACKAGE/data.Version=$VERSION -X $PACKAGE/data.GitCommit=$COMMITHASH"

cYellow='\033[1;33m'
cReset='\033[0m'

log() {
    local msg=$1
    echo -e "${cYellow}${msg}${cReset}"
}

if [[ $1 == "install" ]]; then
    cd "$CWD/cmd/credman"
    log "Compiling credman"
    go install -ldflags="$VERINFO -s -w"
elif [[ $1 == "build" ]]; then
    cd "$CWD/cmd/credman"
    log "Compiling credman"
    go build -ldflags="$VERINFO -s -w"
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

        log "Compiling for $os-$arch"
        GOOS=$os GOARCH=$arch go build -ldflags="$VERINFO -s -w" -o "credman$ext"

        chmod +x "credman$ext"

        log "  Creating archive"
        tar -czf "$filename.tar.gz" "credman$ext"
        mv "$filename.tar.gz" "$CWD/release"

        log "  Removing binary"
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

log "Done"

#!/bin/bash

CWD=$(pwd)

PACKAGE='github.com/ShrewdSpirit/credman'

VERSION='0.10.0'
COMMITHASH=$(git rev-parse HEAD)

VERINFO="-X $PACKAGE/data.Version=$VERSION -X $PACKAGE/data.GitCommit=$COMMITHASH"

cYellow='\033[1;33m'
cReset='\033[0m'

log() {
    local msg=$1
    echo -e "${cYellow}${msg}${cReset}"
}

cd gui/assets
log "Building GUI"
npm run build
log "Embedding assets"
gassets -d .

if [[ $# -eq 1 && $1 == "install" ]]; then
    cd "$CWD/cmd/credman"
    log "Compiling credman"
    go install -ldflags="$VERINFO" -tags='gui'
else
    if [ -d "release" ]; then
        rm -fr release
    fi

    mkdir release

    build() {
        local os=$1
        local arch=$2
        local tags=$3
        local filename="credman-$os-$arch"
        local ext=""

        cd "$CWD/cmd/credman"

        if [[ $os == "windows" ]]; then
            ext=".exe"
        fi

        log "Compiling for $os-$arch"
        GOOS=$os GOARCH=$arch go build -ldflags="$VERINFO -s -w" -o "credman$ext" -tags=$tags

        chmod +x "credman$ext"

        # cd "$CWD/release"
        log "  Creating archive"
        tar -czf "$filename.tar.gz" "credman$ext"
        mv "$filename.tar.gz" "$CWD/release"

        log "  Removing binary"
        rm -fr "credman$ext" "$filename.tar.gz"
    }

    guiTag='gui'
    noGui=''

    build linux amd64 $guiTag
    build linux 386 $guiTag
    build linux arm $noGui
    build linux arm64 $noGui
    build windows amd64 $guiTag
    build windows 386 $guiTag
    build darwin amd64 $noGui
    build darwin 386 $noGui
fi

log "Done"

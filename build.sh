#!/bin/bash

cd app

for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        BIN_FILENAME="itupod-${GOOS}-${GOARCH}"
        if [[ "${GOOS}" == "windows" ]]; then BIN_FILENAME="${BIN_FILENAME}.exe"; fi
        GOOS=${GOOS} GOARCH=${GOARCH} go build -v -o ../bin/$BIN_FILENAME
    done
done

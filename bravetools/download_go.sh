#!/bin/bash

VERSION="1.18.1"
TARGETARCH="arm64"

ARCH=$(uname -m)
if [ $ARCH != "aarch64" ]
then
    TARGETARCH="amd64"
fi

wget -L "https://golang.org/dl/go${VERSION}.linux-${TARGETARCH}.tar.gz"
tar -xf "go${VERSION}.linux-${TARGETARCH}.tar.gz"
rm go${VERSION}.linux-${TARGETARCH}.tar.gz

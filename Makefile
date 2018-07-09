PKG=github.com/meeDamian/bc1toolkit

VERSION_STAMP="${PKG}/lib/common.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="${PKG}/lib/common.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="${PKG}/lib/common.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"

bin/bc1tx: bc1tx/main.go lib/common/common.go
	go build -v -o $@ -ldflags ${BUILD_FLAGS} "${PKG}/bc1tx"

bin/bc1isup: bc1isup/main.go lib/common/common.go
	go build -v -o $@ -ldflags ${BUILD_FLAGS} "${PKG}/bc1isup"

all: bin/bc1tx bin/bc1isup

.PHONY: all

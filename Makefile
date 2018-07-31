PKG=github.com/meeDamian/bc1toolkit

VERSION_STAMP="${PKG}/lib/help.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="${PKG}/lib/help.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="${PKG}/lib/help.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"

SRC_LIB := $(shell find lib -type f -name '*.go')

# currently supported platforms
platforms = windows-amd64.exe darwin-amd64 linux-amd64 linux-arm freebsd-amd64


#
## Simple builds for the current platform
#
bin/bc1isup: bc1isup/main.go $(SRC_LIB)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

all: bin/bc1isup


#
## build for all platforms for a given binary
#
bc1isup = $(addprefix release/bc1isup-,$(platforms))
$(bc1isup): bc1isup/main.go $(SRC_LIB)
	env GOARCH=$(subst .exe,,$(lastword $(subst -, ,$@))) \
	GOOS=$(lastword $(filter-out $(lastword $(subst -, ,$@)), $(subst -, ,$@))) \
	GOARM=$(subst arm,5,$(filter arm,$(subst .exe,,$(lastword $(subst -, ,$@))))) \
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

release/bc1isup: is-git-clean $(bc1isup)

is-git-clean:
	git diff-index --quiet HEAD


dist: release/bc1isup
	zip release/bc1toolkit-mac.zip $(wildcard release/*-darwin-amd64)
	zip release/bc1toolkit-linux.zip $(wildcard release/*-linux-amd64)
	zip release/bc1toolkit-raspberry.zip $(wildcard release/*-linux-arm)
	zip release/bc1toolkit-windows.zip $(wildcard release/*-windows-amd64.exe)
	zip release/bc1toolkit-freebsd.zip $(wildcard release/*-freebsd-amd64)


clean:
	rm -f bin/*
	rm -f release/*


install:
	go install -v -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

# TODO: uninstall:

.PHONY: all is-git-clean dist clean install \
		release/bc1isup \
		release/bc1toolkit-mac.zip




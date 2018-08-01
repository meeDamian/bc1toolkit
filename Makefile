PKG=github.com/meeDamian/bc1toolkit

VERSION_STAMP="${PKG}/lib/help.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="${PKG}/lib/help.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="${PKG}/lib/help.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"

SRC_LIB := $(shell find lib -type f -name '*.go')
VENDOR_LIB := $(shell find vendor -type f -name '*.go')
ALL_SRC := $(shell find . -type f -name '*.go')

# currently supported platforms
platforms = windows-amd64.exe darwin-amd64 linux-amd64 linux-arm freebsd-amd64
binaries = bc1isup bc1tx


#
## Simple builds for the current platform
#
# TODO: collapse all into one
bin/bc1isup: bc1isup/main.go $(SRC_LIB) $(VENDOR_LIB)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

bin/bc1tx: bc1tx/main.go $(SRC_LIB) $(VENDOR_LIB)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1tx

all: bin/bc1isup bin/bc1tx


#
## build for all platforms for each binary
#
# combine all binaries with all platforms they'll be built for
allTargets = $(foreach binary,$(addprefix release/,$(binaries)),$(addprefix $(binary)-,$(platforms)))
$(allTargets): is-git-clean $(ALL_SRC)
	env GOARCH=$(subst .exe,,$(lastword $(subst -, ,$@))) \
	GOOS=$(lastword $(filter-out $(lastword $(subst -, ,$@)), $(subst -, ,$@))) \
	GOARM=$(subst arm,5,$(filter arm,$(subst .exe,,$(lastword $(subst -, ,$@))))) \
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/$(firstword $(subst -, ,$(subst release/,,$@)))

is-git-clean:
	git diff-index --quiet HEAD


dist: $(allTargets)
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

# TODO: uninstall target

.PHONY: all is-git-clean dist clean install




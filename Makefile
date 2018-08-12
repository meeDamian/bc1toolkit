PKG=github.com/meeDamian/bc1toolkit

VERSION_STAMP="${PKG}/lib/help.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="${PKG}/lib/help.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="${PKG}/lib/help.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"

SRC_LIB := $(shell find lib -type f -name '*.go')
VENDOR_LIB := $(shell find vendor -type f -name '*.go')
ALL_SRC := $(shell find . -type f -name '*.go')

# currently supported platforms
platforms = windows-amd64 darwin-amd64 linux-amd64 linux-arm freebsd-amd64
binaries = bc1isup


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
tmpTargets = $(addprefix dist/,$(foreach platform,$(platforms),$(addprefix $(platform)/,$(binaries))))

# add `.exe` to windows binary. Thanks, Windows.
allTargets = $(filter-out dist/windows%,$(tmpTargets)) $(addsuffix .exe,$(filter dist/windows%,$(tmpTargets)))
$(allTargets): $(ALL_SRC)
	env GOARCH=$(word 3,$(subst /, ,$(subst -, ,$@))) \
	GOOS=$(word 2,$(subst /, ,$(subst -, ,$@))) \
	GOARM=$(subst arm,5,$(filter arm,$(word 3,$(subst /, ,$(subst -, ,$@))))) \
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/$(lastword $(subst .exe,,$(subst /, ,$@)))


is-git-clean:
	git diff-index --quiet HEAD


dist: clean $(allTargets)
	zip -j dist/bc1toolkit-mac.zip dist/darwin-amd64/*
	zip -j dist/bc1toolkit-linux.zip dist/linux-amd64/*
	zip -j dist/bc1toolkit-raspberry.zip dist/linux-arm/*
	zip -j dist/bc1toolkit-windows.zip dist/windows-amd64/*
	zip -j dist/bc1toolkit-freebsd.zip dist/freebsd-amd64/*


clean:
	rm -f bin/*
	rm -rf dist/*


install:
	go install -v -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

# TODO: uninstall target

.PHONY: all is-git-clean dist clean install




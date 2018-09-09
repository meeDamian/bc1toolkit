PKG=github.com/meeDamian/bc1toolkit

##
## IMPORTANT:
##		INCREMENT THIS PRIOR TO A NEW RELEASE
##
VERSION := v0.0.3

VERSION_VERSION="${PKG}/lib/help.version=${VERSION}"
VERSION_STAMP="${PKG}/lib/help.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="${PKG}/lib/help.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="${PKG}/lib/help.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-s -w -X ${VERSION_VERSION} -X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"

SRC_LIB := $(shell find lib -type f -name '*.go')
ALL_SRC := $(shell find . -type f -name '*.go')
EXPLORE_TEMPLATES := $(shell find bc1explore/templates -type f -name '*.html')
GO_MOD := go.mod go.sum

# currently supported platforms
platforms = windows-amd64 darwin-amd64 linux-amd64 linux-arm freebsd-amd64
binaries = bc1isup bc1explore

#
## Code Generation
#
bc1explore/templates_generated.go: bc1explore/templates/gen.go $(EXPLORE_TEMPLATES)
	go generate github.com/meeDamian/bc1toolkit/bc1explore


#
## Simple builds for the current platform
#
# TODO: collapse all into one(?)
bin/bc1isup: bc1isup/main.go $(SRC_LIB) $(GO_MOD)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup

bin/bc1tx: bc1tx/main.go $(SRC_LIB) $(GO_MOD)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1tx

bin/bc1explore: bc1explore/main.go bc1explore/templates_generated.go $(SRC_LIB) $(GO_MOD)
	go build -v -o $@ -ldflags ${BUILD_FLAGS} ${PKG}/bc1explore


all: bin/bc1isup bin/bc1explore


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

git-tag:
	git tag -sa ${VERSION} -m "${VERSION}"

git-push-tag:
	git push origin $(VERSION)


releases/$(VERSION):
	[ ! -d $@ ]
	mkdir -p $@

releases/$(VERSION)/bc1toolkit-$(VERSION)-mac.zip:
	zip -j $@ dist/darwin-amd64/*

releases/$(VERSION)/bc1toolkit-$(VERSION)-linux.zip:
	zip -j $@ dist/linux-amd64/*

releases/$(VERSION)/bc1toolkit-$(VERSION)-raspberry.zip:
	zip -j $@ dist/linux-arm/*

releases/$(VERSION)/bc1toolkit-$(VERSION)-windows.zip:
	zip -j $@ dist/windows-amd64/*

releases/$(VERSION)/bc1toolkit-$(VERSION)-freebsd.zip:
	zip -j $@ dist/freebsd-amd64/*


releases = $(addsuffix .zip,$(addprefix releases/$(VERSION)/bc1toolkit-$(VERSION)-,mac linux raspberry windows freebsd))
dist: is-git-clean git-tag releases/$(VERSION) clean $(allTargets) $(releases) git-push-tag
	@echo "\n\trelease ${VERSION} complete and files are in $</*\n"


clean:
	rm -f bin/*
	rm -rf dist/*


install:
	go install -v -ldflags ${BUILD_FLAGS} ${PKG}/bc1isup
	go install -v -ldflags ${BUILD_FLAGS} ${PKG}/bc1explore

# TODO: uninstall target

.PHONY: all is-git-clean git-tag git-push-tag dist clean install




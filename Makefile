
VERSION_STAMP="main.buildStamp=$$(date -u '+%Y-%m-%d_%I:%M:%S%p')"
VERSION_HASH="main.gitHash=$$(git rev-parse HEAD)"
VERSION_BUILDER="main.builder=$$(whoami)@$$(hostname)"

BUILD_FLAGS="-X ${VERSION_STAMP} -X ${VERSION_HASH} -X ${VERSION_BUILDER}"


bin/bc1tx: bc1tx/main.go
	go build -v -o $@ -ldflags ${BUILD_FLAGS} github.com/meeDamian/bc1tools/bc1tx
	

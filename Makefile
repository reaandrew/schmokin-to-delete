USERNAME="reaandrew"
PROJECT="surge"
GITHUB_TOKEN=$$GITHUB_TOKEN
VERSION=`cat VERSION`
BUILD_TIME=`date +%FT%T%z`
COMMIT_HASH=`git rev-parse HEAD`
DIST_NAME_CONVENTION="dist/{{.OS}}_{{.Arch}}_{{.Dir}}"

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
SOURCES += VERSION
# This is how we want to name the binary output
BINARY=${PROJECT}

# These are the values we want to pass for Version and BuildTime

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-X main.CommitHash=${COMMIT_HASH} -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): deps $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY}

.PHONY: build_quick
build_quick: 
	go build

.PHONY: proto
proto:
	 cd server && protoc --go_out=plugins=grpc:. *.proto

.PHONY: install_linter
install_linter:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

.PHONY: lint
lint:
	golangci-lint run

surge:
	go build

.PHONY: test
test: build_quick
	SURGE_PATH="${CURDIR}/surge" go test -v ./...

APP = kafcat
FOLDERS = ./cmd/... ./pkg/...

BUILD_DIR ?= build

# get version info from git's tags
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_TAG := $(shell git describe --tags --dirty 2>/dev/null)
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null)

# inject version info into version vars
LD_RELEASE_FLAGS += -X github.com/rollulus/kafcat/cmd.GitCommit=${GIT_COMMIT}
LD_RELEASE_FLAGS += -X github.com/rollulus/kafcat/cmd.GitTag=${GIT_TAG}
LD_RELEASE_FLAGS += -X github.com/rollulus/kafcat/cmd.SemVer=${VERSION}

.PHONY: default bootstrap clean test build static docker

default: build

bootstrap:
	glide install -v

clean:
	rm -rf $(BUILD_DIR)

test:
	go test -cover $(FOLDERS)

build:
	go build -ldflags "$(LD_RELEASE_FLAGS)" -o ./$(BUILD_DIR)/$(APP)

publish_github:
	go get github.com/goreleaser/goreleaser
	./goreleaser.yaml.sh "$(LD_RELEASE_FLAGS)" >/tmp/gorel.yaml
	goreleaser --config /tmp/gorel.yaml

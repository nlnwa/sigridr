PROJ=sigridr
ORG_PATH=github.com/nlnwa
REPO_PATH=$(ORG_PATH)/$(PROJ)
VERSION ?= $(shell ./scripts/git-version)

LD_FLAGS="-w -X $(REPO_PATH)/version.Version=$(VERSION)"

.PHONY: release-binary
release-binary:
	@go get github.com/golang/dep/cmd/dep
	@dep ensure
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/...

PROJ:=sigridr
ORG_PATH:=github.com/nlnwa
REPO_PATH:=$(ORG_PATH)/$(PROJ)
VERSION?=$(shell ./scripts/git-version)

## https://golang.org/cmd/link/
## -w Omit the DWARF symbol table.
## -X Set the value of the string variable in importpath named name to value.
LD_FLAGS:= "-w -X $(REPO_PATH)/version.Version=$(VERSION)"

.PHONY: release-binary install-dep api clean-api

## -v print the names of packages as they are complied.
install:
	@go install -v -ldflags $(LD_FLAGS) $(REPO_PATH)/cmd/...

api:
	@$(MAKE) -C ./api

install-dep:
	@go get github.com/golang/dep/cmd/dep
	@dep ensure

release-binary: api install-dep install

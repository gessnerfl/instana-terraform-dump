    
export CGO_ENABLED:=0
export GO111MODULE=on
#export GOFLAGS=-mod=vendor

ifdef TRAVIS_TAG
	VERSION=$(TRAVIS_TAG)
else
	VERSION=$(shell git describe --tags --match "v*" --always --dirty)
endif

.PHONY: all
all: build test vet lint fmt

.PHONY: build
build: clean bin/instana-terraform-dump

bin/instana-terraform-dump:
	@echo "+++++++++++  Run GO Build +++++++++++ "
	@go build -o $@ github.com/gessnerfl/instana-terraform-dump

.PHONY: test
test:
	@echo "+++++++++++  Run GO Test +++++++++++ "
	@go test ./... -cover

.PHONY: vet
vet:
	@echo "+++++++++++  Run GO VET +++++++++++ "
	@go vet -all ./...

.PHONY: lint
lint:
	@echo "+++++++++++  Run GO Lint +++++++++++ "
	@golint -set_exit_status `go list ./...`

.PHONY: fmt
fmt:
	@echo "+++++++++++  Run GO FMT +++++++++++ "
	@go fmt ./...

.PHONY: update
update:
	@GOFLAGS="" go get -u
	@go mod tidy

.PHONY: vendor
vendor:
	@go mod vendor

.PHONY: clean
clean:
	@echo "+++++++++++  Clean up project +++++++++++ "
	@rm -rf bin
	@rm -rf output

.PHONY: release
release: \
	clean \
	instana-terraform-dump-linux-amd64 \
	instana-terraform-dump-darwin-amd64 \
	instana-terraform-dump-windows-amd64

instana-terraform-dump-linux-amd64: GOARGS = GOOS=linux GOARCH=amd64
instana-terraform-dump-linux-amd64: FILE_EXTENSION=""
instana-terraform-dump-darwin-amd64: GOARGS = GOOS=darwin GOARCH=amd64
instana-terraform-dump-darwin-amd64: FILE_EXTENSION=""
instana-terraform-dump-windows-amd64: GOARGS = GOOS=windows GOARCH=amd64
instana-terraform-dump-windows-amd64: FILE_EXTENSION=".exe"
instana-terraform-dump-%:
	@echo "+++++++++++ Build Release $@ +++++++++++ "
	$(GOARGS) go build -o output/$@-$(VERSION)$(FILE_EXTENSION) github.com/gessnerfl/instana-terraform-dump
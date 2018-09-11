GOPATH      := $(GOPATH)
NOW         := $(shell date +%s)

.PHONY: run
run: prerequisites
	@go run cmd/playground/main.go

.PHONY: playground
playground:
	@go run cmd/playground/main.go

.PHONY: vendor
vendor:
	GO111MODULE=on go mod init
	GO111MODULE=on go mod vendor

.PHONY: revendor
revendor:
	rm -rf vendor go.*
	$(MAKE) vendor

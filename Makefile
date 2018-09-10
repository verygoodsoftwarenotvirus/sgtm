GOPATH      := $(GOPATH)
NOW         := $(shell date +%s)

.PHONY: build
build:
	@mkdir -p pkg/convo
	protoc --go_out=plugins=grpc:pkg/convo ./protofiles/*.proto # --plugin=protoc-gen-grpc-web=
	@mv pkg/convo/protofiles/*.go pkg/convo
	@rm -rf pkg/convo/protofiles

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

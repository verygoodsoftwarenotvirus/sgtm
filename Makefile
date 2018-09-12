GOPATH      := $(GOPATH)
NOW         := $(shell date +%s)

# demos

.PHONY: hello-world
hello-world:
	@go run cmd/cli/main.go read --file=example_files/hello_world.go

.PHONY: more
more:
	@go run cmd/cli/main.go read --file=example_files/slightly_more_complicated.go

.PHONY: one-part
one-part:
	@go run cmd/cli/main.go read --file=example_files/slightly_more_complicated.go --part Person

.PHONY: one-polly
one-polly:
	@go run cmd/cli/main.go read --file=example_files/slightly_more_complicated.go --part Person --voice-service=polly

.PHONY:the-goal
the-goal:
	@go run cmd/cli/main.go read --file=example_files/grand_goal.go --voice-service=polly

.PHONY: introspect
introspect:
	@go run cmd/cli/main.go read --file=pkg/interpreter/interpreter.go --part Interpreter --voice-service=polly

.PHONY: multi-introspect
multi-introspect:
	@go run cmd/cli/main.go read --file=pkg/interpreter/interpreter.go --part Interpreter --part Describer --voice-service=polly

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

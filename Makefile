GO ?= go

.PHONY: test fmt run-basic run-interactive

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...

run-basic:
	$(GO) run ./examples/basic

run-interactive:
	$(GO) run ./examples/interactive

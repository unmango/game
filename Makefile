BUF    := go tool buf
DEVCTL := go tool devctl

GO_SRC    != $(DEVCTL) list --go
PROTO_SRC != $(BUF) ls-files

build: ${PROTO_SRC}
	$(BUF) build $?

fmt: ${PROTO_SRC}
	$(BUF) format --write $?

lint: ${PROTO_SRC}
	$(BUF) lint $?

tidy: go.mod ${GO_SRC}
	go mod tidy

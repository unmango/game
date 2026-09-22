PROTO_SRC != $(BUF) ls-files
GO_SRC ?= $(shell find . -name '*.go')

build:
	nix build .#
	$(BUF) build $?

test:
	go tool ginkgo run -r

update:
	nix flake update

check lint:
	nix flake check
	buf lint $?

format fmt:
	nix fmt

tidy: go.sum nix/gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

nix/gomod2nix.toml: go.sum ${GO_SRC}
	gomod2nix generate --dir ${CURDIR} --outdir ${@D}

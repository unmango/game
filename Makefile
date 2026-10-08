LIKEC4    ?= npx -y likec4@1.59.4
ARCH_DIR  ?= docs/architecture

PROTO_SRC ?= $(shell buf ls-files)
GO_SRC    ?= $(shell find . -name '*.go')

build:
	nix build .#
	buf build $?

generate gen:
	buf generate

test:
	go tool ginkgo run -r

update:
	nix flake update

check lint:
	nix flake check
	buf lint $?
	$(MAKE) arch-check

format fmt:
	nix fmt
	${LIKEC4} format ${ARCH_DIR}

arch:
	${LIKEC4} start ${ARCH_DIR}

arch-check:
	${LIKEC4} validate ${ARCH_DIR}

arch-build:
	${LIKEC4} build -o ${ARCH_DIR}/dist ${ARCH_DIR}

tidy: go.sum nix/gomod2nix.toml

go.sum: go.mod ${GO_SRC}
	go mod tidy

nix/gomod2nix.toml: go.sum ${GO_SRC}
	gomod2nix generate --dir ${CURDIR} --outdir ${@D}
